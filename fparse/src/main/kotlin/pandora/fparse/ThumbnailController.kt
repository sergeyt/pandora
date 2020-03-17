package pandora.fparse

import org.apache.commons.io.FileUtils
import org.apache.pdfbox.pdmodel.PDDocument
import org.apache.pdfbox.pdmodel.graphics.image.PDImageXObject
import org.apache.pdfbox.rendering.ImageType
import org.apache.pdfbox.rendering.PDFRenderer
import org.springframework.http.MediaType
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RestController
import java.io.ByteArrayOutputStream
import java.io.File
import javax.imageio.ImageIO
import javax.ws.rs.NotSupportedException


data class ThumbnailRequest(val url: String, val format: String)
data class ThumbnailResult(val id: String, val url: String)

// Downloads file from given URL like pre-signed S3 URL
// Parses file content using Apache Content
// Input JSON {url, options?}
// Returns JSON {metadata, text}
@RestController
class ThumbnailController {
    // TODO stream result right to http response
    @PostMapping("/api/tika/thumbnail", consumes = ["application/json"], produces = ["application/json"])
    fun thumbnail(@RequestBody req: ThumbnailRequest): ByteArray {
        val fileRes = downloadFile(req.url)

        val mediaTypes = MediaType.parseMediaTypes(fileRes.headers["Content-Type"])
        if (!mediaTypes.any { it.isCompatibleWith(MediaType.APPLICATION_PDF) }) {
            throw NotSupportedException("only pdf is supported for now")
        }

        val format = if (req.format === "") "JPEG" else req.format

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
        // FileUtils.writeByteArrayToFile(File("/Users/admin/tmp/thumb.jpg"), bytes)
        return bytes
    }
}
