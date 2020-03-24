package pandora.fparse

import com.fasterxml.jackson.annotation.JsonIgnore
import org.apache.pdfbox.pdmodel.PDDocument
import org.apache.pdfbox.pdmodel.graphics.image.PDImageXObject
import org.apache.pdfbox.rendering.ImageType
import org.apache.pdfbox.rendering.PDFRenderer
import org.springframework.core.io.Resource
import org.springframework.http.HttpEntity
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpMethod
import org.springframework.http.MediaType
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RestController
import org.springframework.web.client.RestTemplate
import java.io.ByteArrayOutputStream
import java.io.File
import java.net.URL
import java.nio.file.Files
import java.nio.file.Path
import javax.imageio.ImageIO
import javax.ws.rs.NotSupportedException

val defaultThumbnailFormat = "JPG"

data class ThumbnailRequest(val url: String, val format: String? = defaultThumbnailFormat)
data class ThumbnailResult(val url: String, @JsonIgnore val body: ByteArray)

// Downloads file from given URL like pre-signed S3 URL
// Parses file content using Apache Content
// Input JSON {url, options?}
// Returns JSON {url}
@RestController
class ThumbnailController {
    // TODO don't create thumbnail if it is already exist
    // TODO stream result right to http response
    @PostMapping("/api/tika/thumbnail", consumes = ["application/json"], produces = ["application/json"])
    fun thumbnail(@RequestBody req: ThumbnailRequest): ThumbnailResult {
        val fileRes = downloadFile(req.url)

        if (!fileRes.mediaType.isCompatibleWith(MediaType.APPLICATION_PDF)) {
            throw NotSupportedException("only pdf is supported for now")
        }

        val format = if (req.format.isNullOrBlank()) defaultThumbnailFormat else req.format
        val doc = PDDocument.load(fileRes.body.inputStream)

        try {
            var pageIndex = doc.pages.indexOfFirst {
                val page = it
                page.resources.xObjectNames.any {
                    val xobj = page.resources.getXObject(it)
                    xobj is PDImageXObject
                }
            }
            if (pageIndex < 0) {
                pageIndex = 0
            }

            val pr = PDFRenderer(doc)
            val bi = pr.renderImageWithDPI(pageIndex, 300F, ImageType.ARGB)

            val out = ByteArrayOutputStream()
            ImageIO.write(bi, format, out)
            out.flush()

            val bytes = out.toByteArray()

            val thumbUrl = saveThumbnail(req, bytes, format, fileRes.name)

            return ThumbnailResult(thumbUrl, bytes)
        } finally {
            doc.close()
        }
    }

    private fun saveThumbnail(info: ThumbnailRequest, body: ByteArray, format: String, fileName: String): String {
        val rest = RestTemplate()
        val headers = HttpHeaders()
        headers.set("Accept", MediaType.APPLICATION_JSON_VALUE)
        headers.set("Content-Type", imageMediaType(format))
        headers.set("Authorization", """Bearer ${systemToken()}""")
        headers.set("X-API-Key", getApiKey())

        val req = HttpEntity(body, headers)
        val thumbUrl = thumbnailUrl(info.url, fileName)
        val resp = rest.exchange(thumbUrl, HttpMethod.POST, req, Resource::class.java);
        if (resp.body == null) {
            throw NullPointerException("expect body")
        }

        return thumbUrl
    }

    private fun imageMediaType(format: String): String {
        return Files.probeContentType("""file.${format.toLowerCase()}""".toPath())
    }

    private fun thumbnailUrl(fileUrl: String, fileName: String): String {
        val f = URL(toAbsoluteUrl(fileUrl))
        var path = f.path
        val prefix = "/api/file"
        if (path.startsWith(prefix)) {
            path = path.substring(prefix.length)
        }
        if (path.matches(Regex("""^/0x\d+$"""))) {
            path = "/$fileName"
        }
        return fileServiceHost() + "/api/file/thumbnails" + path
    }
}

fun String.toPath(): Path {
    return File(this).toPath()
}

data class LoginResponse(val token: String)

fun systemToken(): String {
    val rest = RestTemplate()
    val headers = HttpHeaders()
    headers.set("Accept", MediaType.APPLICATION_JSON_VALUE)
    headers.setBasicAuth("system", System.getenv("SYSTEM_PWD"))
    headers.set("X-API-Key", getApiKey())

    val req = HttpEntity("", headers)
    val loginUrl = fileServiceHost() + "/api/login"
    val resp = rest.exchange(loginUrl, HttpMethod.POST, req, LoginResponse::class.java);
    if (resp.body == null) {
        throw NullPointerException("expect body")
    }

    return resp.body!!.token
}

fun getApiKey(): String {
    return System.getenv("API_KEY")
}
