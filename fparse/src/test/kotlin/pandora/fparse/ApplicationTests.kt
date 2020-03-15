package pandora.fparse

import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest
import java.io.File
import java.nio.file.Files

val aliceUrl = "https://www.adobe.com/be_en/active-use/pdf/Alice_in_Wonderland.pdf"

@SpringBootTest
class ApplicationTests {
    @Test
    fun parseAlice() {
        val ctrl = ParseController()
        val result = ctrl.parse(aliceUrl)
        val creator = result.metadata["creator"]
        assert(creator != null)
        assert(creator is String)
        assert((creator as String).startsWith("Lewis Carroll"))
        assert(result.text.strip().startsWith("BY LEWIS CARROLL"))
    }

    // @Test
    fun parseBooks() {
        val dir = "/Users/admin/Dropbox/books"
        val keys = HashSet<String>()
        Files.list(File(dir).toPath()).forEach {
            val ctrl = ParseController()
            val result = ctrl.parse("file://" + it)
            keys.addAll(result.metadata.keys)
            result.metadata.forEach {
                println("%s=%s".format(it.key, it.value))
            }
            println(result.metadata)
        }
        keys.forEach {
            println(it)
        }
    }

    @Test
    fun generateThumbnail() {
        val ctrl = ThumbnailController()
        val result = ctrl.thumbnail(ThumbnailRequest(aliceUrl, "JPEG"))
        assert(result != null)
    }
}
