import static org.junit.Assert.*;
import java.io.IOException;
import org.junit.Test;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.JsonParser;

import io.sysl.issue66.Article;
import io.sysl.issue66.ArticleModelJsonDeserializer;
import io.sysl.issue66.Content;

public class TestArticle
{
    @Test
    public void testJson() throws IOException {
        Article a1 = buildArticle();
        Article a2 = articleFromJson();
        assertEquals(a1, a2);
    }

    public static Article articleFromJson() throws IOException {
        JsonFactory factory = new JsonFactory();
        String s = "{\"content\" : { \"Text\" : \"Peace and love\" } }";
        JsonParser p = factory.createParser(s);
        p.nextToken();
        ArticleModelJsonDeserializer ad = new ArticleModelJsonDeserializer();
        return ad.deserialize(p, (Article)null);
    }

    public static Article buildArticle() {
        Content content = new Content();
        content.setText("Peace and love");
        Article a = new Article();
        a.setContent(content);
        return a;
    }
}
