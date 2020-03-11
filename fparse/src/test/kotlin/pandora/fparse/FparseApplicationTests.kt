package pandora.fparse

import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest

@SpringBootTest
class FparseApplicationTests {
    @Test
    fun parseAlice() {
        val url = "https://www.adobe.com/be_en/active-use/pdf/Alice_in_Wonderland.pdf"
        val ctrl = FparseController()
        val result = ctrl.parse(url)
        assert(result.metadata["creator"]!!.first().startsWith("Lewis Carroll"))
        assert(result.text.strip().startsWith("BY LEWIS CARROLL"))
    }
}
