
import java.io.IOException;
import java.io.StringWriter;
import java.math.BigDecimal;
import java.io.File;

import java.nio.file.Files;
import java.nio.file.Paths;

import com.fasterxml.jackson.core.*;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.test.example.tuple.complex.BuyItemFormComplex;
import com.test.example.tuple.complex.UserFormComplex;
import com.test.example.tuple.complex.CustomerType;
import com.test.example.tuple.complex.Email;
import com.test.example.tuple.complex.UserFormComplexJsonDeserializer;
import com.test.example.tuple.complex.UserFormComplexJsonSerializer;

import org.joda.time.LocalDate;

public class Test
{
    public static void main(String[] args) {
        try {
            BuyItemFormComplex e1 = loadBuyItemFormComplex("src/main/resources/test.json");
            BuyItemFormComplex e2 = buildOne();

            String s1 = serialize(e1);
            String s2 = serialize(e2);
            System.out.println(s1);
            System.out.println(s2);
            System.out.println(s1.equals(s2));

            BuyItemFormComplex e3 = loadBuyItemFormComplex("src/main/resources/test_2.json");
            System.out.println(serialize(e3));

        } catch (IOException e) {
            System.out.println(e);
        }
    }

    public static String serialize(BuyItemFormComplex entity) throws IOException {
        JsonFactory factory = new JsonFactory();
        StringWriter w = new StringWriter();
        JsonGenerator generator = factory.createGenerator(w);
        generator.useDefaultPrettyPrinter();
        UserFormComplexJsonSerializer ser = new UserFormComplexJsonSerializer();
        ser.serialize(generator, entity);
        generator.close();
        return w.toString();
    }

    public static BuyItemFormComplex buildOne() {
        Email one = new Email();
        one.setEmail(new String("john.smith@anz.com"));
        Email two = new Email();
        two.setEmail(new String("jsmith@anz.com"));
        Email.Set emailSet = new Email.Set();
        emailSet.add(one);
        emailSet.add(two);

        BuyItemFormComplex e = new BuyItemFormComplex();
        e.setAmount(new BigDecimal(10.11))
            .setFirstName( new String("John"))
            .setLastName(new String("Smith"))
            .setEmails(emailSet)
            .setCustomerType( CustomerType.from(0))
            .setDateOfBirth(LocalDate.parse("2000-02-29"));
        return e;
    }

    public static BuyItemFormComplex loadBuyItemFormComplex(String filename) throws IOException {
        // ObjectMapper mapper = new ObjectMapper();
        // SimpleModule module = new SimpleModule();
        // module.addDeserializer(UserFormComplex.class, new UserFormComplexJsonDeserializer());
        // mapper.registerModule(module);
        // return mapper.readValue(json, BuyItemFormComplex.class);
        JsonFactory factory = new JsonFactory();
        JsonParser p = factory.createParser(new File(filename));
        p.nextToken();        
        UserFormComplexJsonDeserializer ser = new UserFormComplexJsonDeserializer();
        return ser.deserialize(p, (BuyItemFormComplex)null);
    }
}