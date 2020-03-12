package pandora.fparse

import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest

@SpringBootTest
class FileProcApplicationTests {
    @Test
    fun parseAlice() {
        val url = "https://www.adobe.com/be_en/active-use/pdf/Alice_in_Wonderland.pdf"
        val ctrl = FileProcessingApiController()
        val result = ctrl.parse(url)
        val creator = result.metadata["creator"]
        assert(creator != null)
        assert(creator is String)
        assert((creator as String).startsWith("Lewis Carroll"))
        assert(result.text.strip().startsWith("BY LEWIS CARROLL"))
    }
}
