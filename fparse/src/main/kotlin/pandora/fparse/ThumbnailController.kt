package pandora.fparse

import org.apache.pdfbox.pdmodel.PDDocument
import org.apache.pdfbox.rendering.PDFRenderer
import org.springframework.http.MediaType
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RestController
import java.io.ByteArrayOutputStream
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

        val doc: PDDocument = PDDocument.load(fileRes.body!!.inputStream)

        val format = if (req.format === "") "JPG" else req.format

        // TODO render only first page with image
        val pr = PDFRenderer(doc)
        val bi = pr.renderImageWithDPI(0, 300F)

        val outputStream = ByteArrayOutputStream()
        ImageIO.write(bi, format, outputStream)

        return outputStream.toByteArray()
    }
}
