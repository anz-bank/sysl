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

import org.junit.Test;
import static org.junit.Assert.*;

public class TestJsonSerializer
{
    @Test
    public void testSerializeDeserialize() {
        boolean exception = false;
        try {
            BuyItemFormComplex e1 = loadBuyItemFormComplex("src/main/resources/test.json");
            BuyItemFormComplex e2 = createBuyItemFormComplex();

            String s1 = serialize(e1);
            String s2 = serialize(e2);
            assertTrue(s1.equals(s2));
        } catch (IOException e) {
            exception = true;
        }
        assertFalse(exception);
    }

    @Test
    public void testSubObject() {
        boolean exception = false;
        try {
            BuyItemFormComplex entity = loadBuyItemFormComplex("src/main/resources/test_2.json");
            String str = serialize(entity);
            assertTrue(str.equals(WITH_SUB_OBJECT));

        } catch (IOException e) {
            exception = true;
        }
        assertFalse(exception);
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

    public static BuyItemFormComplex createBuyItemFormComplex() {
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
        JsonFactory factory = new JsonFactory();
        JsonParser p = factory.createParser(new File(filename));
        p.nextToken();
        UserFormComplexJsonDeserializer ser = new UserFormComplexJsonDeserializer();
        return ser.deserialize(p, (BuyItemFormComplex)null);
    }

    public static String WITH_SUB_OBJECT =   "{\n"
                                    + "  \"Amount\" : 10.11,\n"
                                    + "  \"BillingAddress\" : {\n"
                                    + "    \"Address1\" : \"ANZ Building\",\n"
                                    + "    \"Address2\" : \"833 Collins Street\",\n"
                                    + "    \"City\" : \"Melbourne\",\n"
                                    + "    \"Country\" : \"AU\",\n"
                                    + "    \"State\" : \"VIC\",\n"
                                    + "    \"ZipCode\" : \"3000\"\n"
                                    + "  },\n"
                                    + "  \"CustomerType\" : 0,\n"
                                    + "  \"Emails\" : [ {\n"
                                    + "    \"Email\" : \"jsmith@anz.com\"\n"
                                    + "  }, {\n"
                                    + "    \"Email\" : \"john.smith@anz.com\"\n"
                                    + "  } ],\n"
                                    + "  \"FirstName\" : \"John\",\n"
                                    + "  \"LastName\" : \"Smith\"\n"
                                    + "}";
}
