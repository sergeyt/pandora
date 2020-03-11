package pandora.fparse

import org.apache.tika.io.TikaInputStream
import org.apache.tika.metadata.Metadata
import org.apache.tika.parser.AutoDetectParser
import org.apache.tika.parser.ParseContext
import org.apache.tika.sax.BodyContentHandler
import org.springframework.core.io.Resource
import org.springframework.http.HttpEntity
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpMethod
import org.springframework.http.MediaType
import org.springframework.web.bind.annotation.*
import org.springframework.web.client.RestTemplate

data class Request(val url: String)
data class Response(val metadata: Map<String, Array<String>>, val text: String)

// Downloads file from given URL like pre-signed S3 URL
// Parses file content using Apache Content
// Input JSON {url, options?}
// Returns JSON {metadata, text}
@RestController
class FparseController {
    @GetMapping("/api/tika/parse", produces = ["application/json"])
    fun parse(@RequestParam(name="url") url: String): Response {
        return parse(Request(url))
    }

    @PostMapping("/api/tika/parse", consumes = ["application/json"], produces = ["application/json"])
    fun parse(@RequestBody req: Request): Response {
        val rest = RestTemplate()
        val headers = HttpHeaders()
        headers.add("Accept", "*/*")

        val fileReq = HttpEntity("", headers)
        val fileRes = rest.exchange(req.url, HttpMethod.GET, fileReq, Resource::class.java)
        if (fileRes.body == null) {
            throw NullPointerException("expect body")
        }
        val mediaTypes = MediaType.parseMediaTypes(fileRes.headers["Content-Type"])

        val parser = AutoDetectParser()
        val handler = BodyContentHandler(500 * 1024 * 1024)
        val metadata = Metadata()
        metadata.set("Content-Type", mediaTypes[0].type)
        val parseContext = ParseContext()

        val stream = TikaInputStream.get(fileRes.body!!.getInputStream())
        parser.parse(stream, handler, metadata, parseContext)

        // TODO normalize metadata
        // TODO custom converted to json object with all value types preserved
        val meta = metadata.names().map { Pair(it, metadata.getValues(it)) }.toMap()

        return Response(meta, handler.toString())
    }
}
