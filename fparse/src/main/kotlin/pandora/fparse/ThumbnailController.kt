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


data class ThumbnailRequest(val url: String, val format: String)
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

        val mediaTypes = MediaType.parseMediaTypes(fileRes.headers["Content-Type"])
        if (!mediaTypes.any { it.isCompatibleWith(MediaType.APPLICATION_PDF) }) {
            throw NotSupportedException("only pdf is supported for now")
        }

        val format = if (req.format === "") "JPG" else req.format

        val doc = PDDocument.load(fileRes.body!!.inputStream)

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

        val thumbUrl = saveThumbnail(req, bytes, format)

        return ThumbnailResult(thumbUrl, bytes)
    }

    private fun saveThumbnail(info: ThumbnailRequest, body: ByteArray, format: String): String {
        val rest = RestTemplate()
        val headers = HttpHeaders()
        headers.add("Accept", "*/*")
        headers.add("Content-Type", imageMediaType(format))

        val req = HttpEntity(body, headers)
        val thumbUrl = thumbnailUrl(info.url)
        val resp = rest.exchange(thumbUrl, HttpMethod.POST, req, Resource::class.java);
        if (resp.body == null) {
            throw NullPointerException("expect body")
        }

        return thumbUrl
    }

    private fun imageMediaType(format: String): String {
        return Files.probeContentType("""file.${format.toLowerCase()}""".toPath())
    }

    private fun thumbnailUrl(fileUrl: String): String {
        val f = URL(toAbsoluteUrl(fileUrl))
        return fileServiceBaseURL() + "/file/thumbnails" + f.path
    }
}

fun String.toPath(): Path {
    return File(this).toPath()
}
