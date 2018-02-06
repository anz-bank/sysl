import static org.junit.Assert.*;

import java.io.File;
import java.io.IOException;
import java.io.StringWriter;
import java.math.BigDecimal;
import java.util.HashSet;

import com.fasterxml.jackson.core.*;

import org.junit.Test;

import org.joda.time.LocalDate;

import io.sysl.reljam.gen.tuple.complex.BuyItemFormComplex;
import io.sysl.reljam.gen.tuple.complex.CustomerType;
import io.sysl.reljam.gen.tuple.complex.Email;
import io.sysl.reljam.gen.tuple.complex.Address;
import io.sysl.reljam.gen.tuple.complex.UserFormComplexJsonDeserializer;
import io.sysl.reljam.gen.tuple.complex.UserFormComplexJsonSerializer;

public class TestJsonSerializer
{
    @Test
    public void testJson1() throws IOException {
        BuyItemFormComplex e1 = loadBuyItemFormComplex("src/test/resources/test1.json");
        BuyItemFormComplex e2 = createBuyItemFormComplex1();
        assertEquals(e1, e2);

        String s1 = serialize(e1);
        String s2 = serialize(e2);
        assertEquals(s1, s2);
    }

    @Test
    public void testJson2()  throws IOException {
        BuyItemFormComplex e1 = loadBuyItemFormComplex("src/test/resources/test2.json");
        BuyItemFormComplex e2 = createBuyItemFormComplex2();
        assertEquals(e1, e2);

        String s1 = serialize(e1);
        String s2 = serialize(e2);
        assertEquals(s1, s2);
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

    public static BuyItemFormComplex createBuyItemFormComplex1() {
        Email one = new Email();
        one.setEmail(new String("john.smith@anz.com"));
        Email two = new Email();
        two.setEmail(new String("jsmith@anz.com"));
        Email.Set emailSet = new Email.Set();
        emailSet.add(one);
        emailSet.add(two);

        HashSet<String> s = new HashSet<String>();
        s.add("one");
        s.add("two");

        BuyItemFormComplex e = new BuyItemFormComplex();
        e.setAmount(new BigDecimal("10.11"))
            .setFirstName( new String("John"))
            .setLastName(new String("Smith"))
            .setEmails(emailSet)
            .setCustomerType( CustomerType.from(0))
            .setDateOfBirth(LocalDate.parse("2000-02-29"))
            .setTags(s);

        return e;
    }

    public static BuyItemFormComplex createBuyItemFormComplex2() {
        BuyItemFormComplex e = createBuyItemFormComplex1();
        e.setDateOfBirth(null);
        e.setTags(null);

        Address a = new Address();
        a.setAddress1(new String("ANZ Building"));
        a.setAddress2(new String("833 Collins Street"));
        a.setCity(new String("Melbourne"));
        a.setState(new String("VIC"));
        a.setZipCode(new String("3000"));
        a.setCountry(new String("AU"));

        e.setBillingAddress(a);

        return e;
    }

    public static BuyItemFormComplex loadBuyItemFormComplex(String filename) throws IOException {
        JsonFactory factory = new JsonFactory();
        JsonParser p = factory.createParser(new File(filename));
        p.nextToken();
        UserFormComplexJsonDeserializer ser = new UserFormComplexJsonDeserializer();
        return ser.deserialize(p, (BuyItemFormComplex)null);
    }

}
