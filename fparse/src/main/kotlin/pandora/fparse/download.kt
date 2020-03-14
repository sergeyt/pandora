package pandora.fparse

import org.springframework.core.io.Resource
import org.springframework.http.HttpEntity
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpMethod
import org.springframework.http.ResponseEntity
import org.springframework.web.client.RestTemplate
import java.net.MalformedURLException
import java.net.URL
import org.slf4j.LoggerFactory;

val LOGGER = LoggerFactory.getLogger(Application::class.java)

fun downloadFile(url: String): ResponseEntity<Resource> {
    val rest = RestTemplate()
    val headers = HttpHeaders()
    headers.add("Accept", "*/*")

    val req = HttpEntity("", headers)
    val file = toAbsoluteUrl(url)

    LOGGER.info("downloading file {}", file)

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
