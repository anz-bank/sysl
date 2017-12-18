import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.io.StringReader;
import java.io.StringWriter;

import java.math.BigDecimal;

import java.util.Comparator;

import javax.xml.stream.XMLOutputFactory;
import javax.xml.stream.XMLInputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.XMLStreamReader;
import javax.xml.stream.XMLStreamWriter;

import java.math.BigDecimal;

import org.joda.time.LocalDate;

import io.sysl.demo.petshop.facade.*;
import io.sysl.demo.petshop.model.*;
import io.sysl.demo.petshop.views.*;

import io.sysl.demo.petshop.api.PetShop;
import io.sysl.demo.petshop.api.PetShopApiXmlSerializer;

import org.junit.Test;
import static org.junit.Assert.*;

public class TestPetShop {

  private static int fib(int i) {
    if (i < 2)
      return i;
    return fib(i - 1) + fib(i - 2);
  }

  @Test
  public void testPetShopModelSerializer() throws XMLStreamException {
    PetShopModel model = new PetShopModel();
    PetShopFacade facade = new PetShopFacade(model);

    Employee mark = newEmployee(facade, "Mark", 1995, 03, 04);

    Breed labrador = newBreed(facade, "Labrador", "dog", 4);

    labrador
      .setAvgLifespan(new BigDecimal(10.2))
      .setAvgWeight(new BigDecimal(101.1)); // tests rounding down.

    Pet ralph = newPet(facade, "Ralph", labrador, 2014, 02, 11, 4);
    EmployeeTendsPet ecfp_mr = newTends(facade, mark, ralph);

    String model_xml = toXml(model);
    String expected_xml = "<PetShopModel><BreedList><Breed><breedId>2</breedId>"
            + "<avgLifespan>10.20</avgLifespan><avgWeight>101</avgWeight><breedName>Labrador</breedName>"
            + "<numLegs>4</numLegs><species>dog</species></Breed></BreedList><EmployeeList>"
            + "<Employee><employeeId>1</employeeId><dob>1995-03-04</dob><name>Mark</name></Employee></EmployeeList>"
            + "<EmployeeTendsPetList><EmployeeTendsPet><employeeId>1</employeeId><petId>3</petId></EmployeeTendsPet>"
            + "</EmployeeTendsPetList><PetList><Pet><petId>3</petId><breedId>2</breedId><dob>2014-02-11</dob>"
            + "<name>Ralph</name><numLegs>4</numLegs></Pet></PetList></PetShopModel>";

    assertEquals(expected_xml, model_xml);

    PetShopModel model_2 = new PetShopModel();
    fromXml(model_2, model_xml);

    model_2.validate();
    assertEquals(true, model_2.getBreedTable().contains(labrador));

  }

  @Test
  public void testPetShop() throws XMLStreamException {
    PetShopModel model = new PetShopModel();
    PetShopFacade facade = new PetShopFacade(model);

    Employee anne = newEmployee(facade, "Anne", 1993, 11, 20);
    Employee mark = newEmployee(facade, "Mark", 1995, 03, 04);

    Breed labrador = newBreed(facade, "Labrador", "dog", 4);
    Breed doberman = newBreed(facade, "Doberman Pinscher", "dog", 4);
    Breed python = newBreed(facade, "Python", "snake", 0);
    Breed taipan = newBreed(facade, "Taipan", "snake", 0);

    Pet ralph = newPet(facade, "Ralph", labrador, 2014, 02, 11, 4);
    Pet marcy = newPet(facade, "Marcy", labrador, 2014, 02, 11, 3);
    Pet boris = newPet(facade, "Boris", doberman, 2015, 06, 01, 4);
    Pet guido = newPet(facade, "Guido", python, 1956, 01, 31, 0);
    Pet tracy = newPet(facade, "Tracy", taipan, 2009, 05, 05, 0);

    EmployeeTendsPet ecfp_ab = newTends(facade, anne, boris);
    EmployeeTendsPet ecfp_ag = newTends(facade, anne, guido);
    EmployeeTendsPet ecfp_mr = newTends(facade, mark, ralph);
    EmployeeTendsPet ecfp_mm = newTends(facade, mark, marcy);
    EmployeeTendsPet ecfp_mb = newTends(facade, mark, boris);

    assertEquals(
      "Employee(employeeId = 1, dob = 1993-11-20, name = Anne)",
      anne.toString());

    assertEquals(
      "[\n" +
      "  Employee(employeeId = 1, dob = 1993-11-20, name = Anne),\n" +
      "  Employee(employeeId = 2, dob = 1995-03-04, name = Mark)\n" +
      "]",
      model.getEmployeeTable().toString());

    PetShopModelToApi view = new PetShopModelToApi() {
      @Override
      public Integer fibonacci(Integer i) {
        return fib(i);
      }
    };
    PetShop petshop = view.modelToApi(model);
    String ps_xml = toXml(petshop);

    PetShopModel model2 = new PetShopModel();
    fromXml(model2, toXml(model));

    // Index numbers are non-deterministic, so we replace them with ?'s.
    assertEquals("employees:\n" + "- [?] Anne 1993-11-20\n" + "- [?] Mark 1995-03-04\n" + "breeds:\n"
        + "- [?] Doberman Pinscher (dog)\n" + "  pets:\n" + "  - Boris 2015-06-01 (legs: 4; rank: 2)\n"
        + "- [?] Labrador (dog)\n" + "  pets:\n" + "  - Marcy 2014-02-11 (legs: 3; rank: 1)\n"
        + "  - Ralph 2014-02-11 (legs: 4; rank: 2)\n" + "- [?] Python (snake)\n" + "  pets:\n"
        + "  - Guido 1956-01-31 (legs: 0; rank: 0)\n" + "- [?] Taipan (snake)\n" + "  pets:\n"
        + "  - Tracy 2009-05-05 (legs: 0; rank: 0)\n", report(petshop).replaceAll("\\[\\d+\\] ", "[?] "));

    assertEquals(new Integer(11), petshop.getNumLegs());

    assertEquals(new Integer(1), new Integer(labrador.toPetView().toEmployeeTendsPetView().toEmployeeView().size()));
    assertEquals(new Integer(2), new Integer(doberman.toPetView().toEmployeeTendsPetView().toEmployeeView().size()));

    PetShopApiToModel view2 = new PetShopApiToModel();
    PetShopModel model3 = view2.apiToModel(petshop);
    for (Employee e : model3.getEmployeeTable()) {
      assertEquals(0, e.getError().intValue());
    }
  }

  @Test
  public void testAttributes() throws XMLStreamException {
    io.sysl.demo.petshop.api.Employee employee = new io.sysl.demo.petshop.api.Employee();
    employee.setIndex(1);
    employee.setName("Anne");

    assertEquals("<Employee index=\"1\"><name>Anne</name></Employee>", toXml(employee));

  }

  @Test
  public void testFollowNullForeignKey() throws XMLStreamException {
    PetShopModel model = new PetShopModel();
    PetShopFacade facade = new PetShopFacade(model);

    Breed labrador = newBreed(facade, "Labrador", "dog", 4);
    Breed doberman = newBreed(facade, "Doberman Pinscher", "dog", 4);

    Pet ralph = newPet(facade, "Ralph", labrador, 2014, 02, 11, 4);
    Pet marcy = newPet(facade, "Marcy", labrador, 2014, 02, 11, 3);

    Pet.Table pet = model.getPetTable();
    Breed.View breeds = pet.toBreedView();

    assertEquals(1, breeds.size());

    marcy.setBreed(doberman);
    assertEquals(2, breeds.size());

    marcy.setBreed(null);
    assertEquals(1, breeds.size());

    ralph.setBreed(null);
    assertEquals(0, breeds.size());
  }

  private Employee newEmployee(PetShopFacade facade, String name, int y, int m, int d) {
    LocalDate dob = new LocalDate(y, m, d);
    return facade.getEmployee().insert().setName(name).setDob(dob);
  }

  private Breed newBreed(PetShopFacade facade, String name, String species, int numLegs) {
    return facade.getBreed().insert().setBreedName(name).setSpecies(species).setNumLegs(numLegs);
  }

  private Pet newPet(PetShopFacade facade, String name, Breed breed, int y, int m, int d, int numLegs) {
    LocalDate dob = new LocalDate(y, m, d);
    return facade.getPet().insert().setName(name).setBreed(breed).setDob(dob).setNumLegs(numLegs);
  }

  private EmployeeTendsPet newTends(PetShopFacade facade, Employee employee, Pet pet) {
    return facade.getEmployeeTendsPet().build().withEmployee(employee).withPet(pet).insert();
  }

  private String toXml(PetShopModel model) throws XMLStreamException {
    StringWriter sw = new StringWriter();
    XMLOutputFactory xof = XMLOutputFactory.newFactory();
    XMLStreamWriter xsw = xof.createXMLStreamWriter(sw);

    PetShopModelXmlSerializer xmlOut = new PetShopModelXmlSerializer();
    xmlOut.serialize(model, xsw);
    xsw.close();

    return sw.toString();
  }

  private String toXml(PetShop petshop) throws XMLStreamException {
    StringWriter sw = new StringWriter();
    XMLOutputFactory xof = XMLOutputFactory.newFactory();
    XMLStreamWriter xsw = xof.createXMLStreamWriter(sw);

    PetShopApiXmlSerializer xmlOut = new PetShopApiXmlSerializer();
    xmlOut.serialize(petshop, xsw, "PetShop");
    xsw.close();

    return sw.toString();
  }

  private void fromXml(PetShopModel model, String xml) throws XMLStreamException {
    StringReader sr = new StringReader(xml);
    XMLInputFactory xif = XMLInputFactory.newFactory();
    XMLStreamReader xsr = xif.createXMLStreamReader(sr);

    PetShopModelXmlDeserializer xmlIn = new PetShopModelXmlDeserializer();
    xmlIn.deserialize(model, xsr);
  }

  private String toXml(io.sysl.demo.petshop.api.Employee employee) throws XMLStreamException {
    StringWriter sw = new StringWriter();
    XMLOutputFactory xof = XMLOutputFactory.newFactory();
    XMLStreamWriter xsw = xof.createXMLStreamWriter(sw);

    PetShopApiXmlSerializer xmlOut = new PetShopApiXmlSerializer();
    xmlOut.serialize(employee, xsw, "Employee");
    xsw.close();

    return sw.toString();
  }

  private String report(PetShop petshop) {
    ByteArrayOutputStream baos = new ByteArrayOutputStream();
    PrintStream out = new PrintStream(baos);
    out.print("employees:\n");
    for (io.sysl.demo.petshop.api.Employee e :
        petshop.getEmployees().orderBy(employeeComparator)) {
      out.printf("- [%d] %s %tY-%<tm-%<td\n",
        e.getIndex(), e.getName(), e.getDob().toDate());
    }
    out.print("breeds:\n");
    for (io.sysl.demo.petshop.api.Breed b :
        petshop.getBreeds().orderBy(breedComparator)) {
      out.printf("- [%d] %s (%s)\n", b.getIndex(), b.getName(), b.getSpecies());
      out.printf("  pets:\n");
      for (io.sysl.demo.petshop.api.Pet p :
          b.getPets().orderBy(petComparator)) {
        out.printf("  - %s %tY-%<tm-%<td (legs: %d; rank: %d)\n",
          p.getName(), p.getDob().toDate(), p.getNumLegs(), p.getLegRank());
      }
    }
    return baos.toString();
  }

  private final Comparator<io.sysl.demo.petshop.api.Employee>
    employeeComparator =
    new Comparator<io.sysl.demo.petshop.api.Employee>() {
      @Override
      public int compare(
          io.sysl.demo.petshop.api.Employee a,
          io.sysl.demo.petshop.api.Employee b) {
        int c = a.getName().compareTo(b.getName());
        return c != 0 ? c : a.getDob().compareTo(b.getDob());
      }
    };

  private final Comparator<io.sysl.demo.petshop.api.Breed>
    breedComparator =
    new Comparator<io.sysl.demo.petshop.api.Breed>() {
      @Override
      public int compare(
          io.sysl.demo.petshop.api.Breed a,
          io.sysl.demo.petshop.api.Breed b) {
        int c = a.getName().compareTo(b.getName());
        return c != 0 ? c : a.getSpecies().compareTo(b.getSpecies());
      }
    };

  private final Comparator<io.sysl.demo.petshop.api.Pet>
    petComparator =
    new Comparator<io.sysl.demo.petshop.api.Pet>() {
      @Override
      public int compare(
          io.sysl.demo.petshop.api.Pet a,
          io.sysl.demo.petshop.api.Pet b) {
        int c = a.getName().compareTo(b.getName());
        return c != 0 ? c : a.getDob().compareTo(b.getDob());
      }
    };

}
