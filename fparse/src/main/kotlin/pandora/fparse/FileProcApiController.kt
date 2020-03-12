package pandora.fparse

import org.apache.pdfbox.pdmodel.PDDocument
import org.apache.pdfbox.rendering.PDFRenderer
import org.apache.tika.io.TikaInputStream
import org.apache.tika.metadata.Metadata
import org.apache.tika.parser.AutoDetectParser
import org.apache.tika.parser.ParseContext
import org.apache.tika.sax.BodyContentHandler
import org.springframework.core.io.Resource
import org.springframework.http.*
import org.springframework.web.bind.annotation.*
import org.springframework.web.client.RestTemplate
import java.io.ByteArrayOutputStream
import javax.imageio.ImageIO
import javax.ws.rs.NotSupportedException


data class ParseRequest(val url: String)
data class ParseResult(val metadata: Map<String, Any>, val text: String)

data class ThumbnailRequest(val url: String, val format: String)
data class ThumbnailResult(val id: String, val url: String)

fun normalizeValue(vals: Array<String>): Any {
    if (vals.size == 1) {
        val v = vals[0]
        if (v === "true")
            return true
        if (v === "false")
            return false
        return v
    } else {
        return vals
    }
}

fun downloadFile(url: String): ResponseEntity<Resource> {
    val rest = RestTemplate()
    val headers = HttpHeaders()
    headers.add("Accept", "*/*")

    val req = HttpEntity("", headers)
    val res = rest.exchange(url, HttpMethod.GET, req, Resource::class.java)
    if (res.body == null) {
        throw NullPointerException("expect body")
    }

    return res
}

// Downloads file from given URL like pre-signed S3 URL
// Parses file content using Apache Content
// Input JSON {url, options?}
// Returns JSON {metadata, text}
@RestController
class FileProcApiController {
    @GetMapping("/api/tika/parse", produces = ["application/json"])
    fun parse(@RequestParam(name = "url") url: String): ParseResult {
        return parse(ParseRequest(url))
    }

    @PostMapping("/api/tika/parse", consumes = ["application/json"], produces = ["application/json"])
    fun parse(@RequestBody req: ParseRequest): ParseResult {
        val fileRes = downloadFile(req.url)
        val mediaTypes = MediaType.parseMediaTypes(fileRes.headers["Content-Type"])

        val parser = AutoDetectParser()
        val handler = BodyContentHandler(500 * 1024 * 1024)
        val metadata = Metadata()
        metadata.set("Content-Type", mediaTypes[0].type)
        val parseContext = ParseContext()

        val stream = TikaInputStream.get(fileRes.body!!.inputStream)
        parser.parse(stream, handler, metadata, parseContext)

        // TODO normalize metadata, i.e. convert to standard names
        val meta = metadata.names().map {
            val vals = metadata.getValues(it)
            Pair(it, normalizeValue(vals))
        }.toMap()

        return ParseResult(meta, handler.toString())
    }

    // TODO stream result right to http response
    @PostMapping("/api/tika/thumbnail", consumes = ["application/json"], produces = ["application/json"])
    fun thumbnail(@RequestBody req: ThumbnailRequest): ByteArray {
        val fileRes = downloadFile(req.url)

        val mediaTypes = MediaType.parseMediaTypes(fileRes.headers["Content-Type"])
        if (!mediaTypes.any { it.isCompatibleWith(MediaType.APPLICATION_PDF) }) {
            throw NotSupportedException("only pdf is supported for now")
        }

        val doc: PDDocument = PDDocument.load(fileRes.body!!.inputStream)

        val format = if (req.format === "") "JPEG" else req.format

        // TODO render only first page with image
        val pr = PDFRenderer(doc)
        val bi = pr.renderImageWithDPI(0, 300F)

        val outputStream = ByteArrayOutputStream()
        ImageIO.write(bi, format, outputStream)

        return outputStream.toByteArray()
    }
}
