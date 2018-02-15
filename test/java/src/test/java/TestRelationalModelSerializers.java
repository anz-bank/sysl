import static org.junit.Assert.*;

import java.io.File;
import java.io.FileReader;
import java.io.IOException;

import javax.xml.stream.XMLInputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.XMLStreamReader;

import com.fasterxml.jackson.core.*;
import com.fasterxml.jackson.databind.DeserializationContext;

import org.junit.Test;

import io.sysl.facade.UserFacade;
import io.sysl.model.*;

public class TestRelationalModelSerializers
{
    @Test
    public void test1() throws IOException, XMLStreamException {
        UserModel m1 = loadUserModelFromJsonFile("src/test/resources/test3.json");
        UserModel m2 = loadUserModelFromXmlFile("src/test/resources/test4.xml");
        assertEquals(m1, m2);
    }

    @Test
    public void test2() throws IOException, XMLStreamException {
        UserModel m1 = loadUserModelFromJsonFile("src/test/resources/test3.json");
        UserModel m2 = createUserModel();

        assertEquals(m1, m2);
    }

    public static UserModel createUserModel() {
        UserFacade f = new UserFacade(new UserModel());

        Address a1 = f.getAddress().insert();
        a1.setAddress1(new String("3000 Flinders Street"));
        a1.setCity(new String("Melbourne"));
        a1.setState(new String("VIC"));
        a1.setZipCode(new String("3001"));
        a1.setCountry(new String("AU"));

        Users u = f.getUsers().insert();
        u.setFirstName("John").setLastName("Smith").setPrimaryAddressIDFrom(a1);

        Address a2 = f.getAddress().insert();
        a2.setAddress1(new String("ANZ Building"));
        a2.setAddress2(new String("833 Collins Street"));
        a2.setCity(new String("Melbourne"));
        a2.setState(new String("VIC"));
        a2.setZipCode(new String("3000"));
        a2.setCountry(new String("AU"));


        UserAddressRel ua = f.getUserAddressRel().insert();
        ua.setAddress(a1).setUsers(u).setAddressType("Home");

        UserAddressRel ua2 = f.getUserAddressRel().insert();
        ua2.setAddress(a2).setUsers(u).setAddressType("Office");

        return f.getModel();
    }

    public static UserModel loadUserModelFromJsonFile(String filename) throws IOException {
        JsonFactory factory = new JsonFactory();
        JsonParser p = factory.createParser(new File(filename));
        p.nextToken();
        UserModelJsonDeserializer ser = new UserModelJsonDeserializer();
        return ser.deserialize(p, (DeserializationContext)null);
    }

    public static UserModel loadUserModelFromXmlFile(String filename) throws IOException, XMLStreamException {
        XMLInputFactory xif = XMLInputFactory.newFactory();
        XMLStreamReader xsr = xif.createXMLStreamReader(new FileReader(filename));
        UserModel m = new UserModel();

        UserModelXmlDeserializer ser = new UserModelXmlDeserializer();
        ser.deserialize(m, xsr);

        return m;
    }
}
