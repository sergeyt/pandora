package pandora.fparse

import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest

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

    @Test
    fun generateThumbnail() {
        val ctrl = ThumbnailController()
        val result = ctrl.thumbnail(ThumbnailRequest(aliceUrl, "JPEG"))
        assert(result != null)
    }
}
