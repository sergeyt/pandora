package pandora.fparse

import org.slf4j.LoggerFactory
import org.springframework.core.io.FileUrlResource
import org.springframework.core.io.InputStreamSource
import org.springframework.core.io.Resource
import org.springframework.http.*
import org.springframework.web.client.RestTemplate
import java.io.File
import java.net.MalformedURLException
import java.net.URL
import java.nio.file.Files

val LOGGER = LoggerFactory.getLogger(Application::class.java)

data class FileResponse(val name: String, val mediaType: MediaType, val body: InputStreamSource)

fun downloadFile(url: String): FileResponse {
    val file = toAbsoluteUrl(url)

    LOGGER.info("downloading file {}", file)

    // WARNING for testing purposes only
    // TODO in production allow only specific directories
    if (file.startsWith("file://")) {
        val path = file.removePrefix("file://")
        val contentType = determineContentType(path)
        val headers = HttpHeaders()
        headers["Content-Type"] = contentType
        val resource = FileUrlResource(path)
        return FileResponse(path, MediaType.parseMediaType(contentType), resource)
    }

    val rest = RestTemplate()
    val headers = HttpHeaders()
    headers.set("Accept", "*/*")

    val req = HttpEntity("", headers)

    val res = rest.exchange(file, HttpMethod.GET, req, Resource::class.java)
    if (res.body == null) {
        throw NullPointerException("expect body")
    }

    val mediaTypes = MediaType.parseMediaTypes(res.headers["Content-Type"])
    var mediaType = mediaTypes.firstOrNull()
    if (mediaType == null) {
        mediaType = MediaType.APPLICATION_OCTET_STREAM
    }

    // FIXME parse content-disposition header
    val name = res.body!!.filename ?: ""
    return FileResponse(name, mediaType, res.body!!)
}

fun determineContentType(path: String): String {
    if (path.endsWith(".pdf", true)) {
        return MediaType.APPLICATION_PDF.toString()
    }
    return Files.probeContentType(File(path).toPath())
}

fun toAbsoluteUrl(url: String): String {
    try {
        URL(url)
        return url
    } catch (e: MalformedURLException) {
        return fileServiceHost() + url
    }
}

fun fileServiceHost(): String {
    var base = System.getenv("FS_HOST")
    if (base == "") {
        base = "http://localhost"
    }
    return base
}
