package pandora.fparse

import org.apache.tika.io.TikaInputStream
import org.apache.tika.metadata.Metadata
import org.apache.tika.parser.AutoDetectParser
import org.apache.tika.parser.ParseContext
import org.apache.tika.sax.BodyContentHandler
import org.springframework.web.bind.annotation.*


data class ParseRequest(val url: String)
data class ParseResult(val metadata: Map<String, Any>, val text: String)

// Downloads file from given URL like pre-signed S3 URL
// Parses file content using Apache Tika
// Input JSON {url, options?}
// Returns JSON {metadata, text}
@RestController
class ParseController {
    @GetMapping("/api/tika/parse", produces = ["application/json"])
    fun parse(@RequestParam(name = "url") url: String): ParseResult {
        return parse(ParseRequest(url))
    }

    @PostMapping("/api/tika/parse", consumes = ["application/json"], produces = ["application/json"])
    fun parse(@RequestBody req: ParseRequest): ParseResult {
        val fileRes = downloadFile(req.url)

        val parser = AutoDetectParser()
        val handler = BodyContentHandler(500 * 1024 * 1024)
        val metadata = Metadata()
        metadata.set("Content-Type", fileRes.mediaType.type)
        val parseContext = ParseContext()

        val stream = TikaInputStream.get(fileRes.body.inputStream)
        parser.parse(stream, handler, metadata, parseContext)

        val dups = HashSet<String>()
        val meta = metadata.names().map {
            val vals = metadata.getValues(it)
            // TODO allow fullname in some cases
            val name = stem(it.replace(':', '.').toLowerCase().split('.').last())
            var value = normalizeValue(vals)
            val strval = if (value is Array<*>) value.joinToString(";") else value
            // dedupe dc.creator, creator, etc
            var p = name + strval
            if (dups.contains(p)) {
                value = ""
            } else {
                dups.add(p)
            }
            Pair(name, value)
        }.filter { it.second != "" }.toMap()

        return ParseResult(meta, handler.toString())
    }
}

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

fun stem(name: String): String {
    val k = name
            .replace("_", "")
            .replace("-", "")
    if (k === "keywords") {
        return "keyword"
    }
    if (k === "contenttype") {
        return "content_type"
    }
    if (k == "lastmodified" || k == "lastsavedate") {
        return "modified_at"
    }
    if (k == "created" || k == "creationdate") {
        return "created_at"
    }
    return k
}
