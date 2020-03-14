package pandora.fparse

import org.slf4j.LoggerFactory
import org.springframework.core.io.FileUrlResource
import org.springframework.core.io.Resource
import org.springframework.http.*
import org.springframework.web.client.RestTemplate
import java.io.File
import java.net.MalformedURLException
import java.net.URL
import java.nio.file.Files

val LOGGER = LoggerFactory.getLogger(Application::class.java)

fun downloadFile(url: String): ResponseEntity<Resource> {
    val file = toAbsoluteUrl(url)

    LOGGER.info("downloading file {}", file)

    // WARNING for testing purposes only
    // TODO in production allow only specific directories
    if (file.startsWith("file://")) {
        val path = File(file.removePrefix("file://")).toPath()
        val contentType = Files.probeContentType(path)
        val headers = HttpHeaders()
        headers["Content-Type"] = contentType
        val resource = FileUrlResource(path.toString())
        return ResponseEntity(resource, headers, HttpStatus.OK)
    }

    val rest = RestTemplate()
    val headers = HttpHeaders()
    headers.add("Accept", "*/*")

    val req = HttpEntity("", headers)

    val res = rest.exchange(file, HttpMethod.GET, req, Resource::class.java)
    if (res.body == null) {
        throw NullPointerException("expect body")
    }

    return res
}

fun toAbsoluteUrl(url: String): String {
    try {
        URL(url)
        return url
    } catch (e: MalformedURLException) {
        var base = System.getenv("FS_HOST")
        if (base == "") {
            base = "http://localhost"
        }
        return base + url
    }
}
