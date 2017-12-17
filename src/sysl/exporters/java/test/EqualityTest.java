import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.io.StringReader;
import java.io.StringWriter;

import java.util.Comparator;

import javax.xml.stream.XMLOutputFactory;
import javax.xml.stream.XMLInputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.XMLStreamReader;
import javax.xml.stream.XMLStreamWriter;

import org.joda.time.LocalDate;

import io.sysl.demo.petshop.facade.*;
import io.sysl.demo.petshop.model.*;
import io.sysl.demo.petshop.views.*;

import io.sysl.demo.petshop.api.PetShop;
import io.sysl.demo.petshop.api.PetShopApiXmlSerializer;

import org.junit.Test;
import static org.junit.Assert.*;

public class EqualityTest {

  private static class PetShop {
    public PetShopModel model;
    public PetShopFacade facade;

    public Employee anne;
    public Employee mark;

    public Breed labrador;
    public Breed doberman;
    public Breed python  ;
    public Breed taipan  ;

    public Pet ralph;
    public Pet marcy;
    public Pet boris;
    public Pet guido;
    public Pet tracy;

    public EmployeeTendsPet anneTendsBoris;
    public EmployeeTendsPet anneTendsGuido;
    public EmployeeTendsPet markTendsRalph;
    public EmployeeTendsPet markTendsMarcy;
    public EmployeeTendsPet markTendsBoris;

    public static PetShop base() {
      PetShop ps = new PetShop();

      ps.anne = ps.newEmployee("Anne", 1993, 11, 20);
      ps.mark = ps.newEmployee("Mark", 1995, 03, 04);

      ps.labrador = ps.newBreed("Labrador", "dog", 4);
      ps.doberman = ps.newBreed("Doberman Pinscher", "dog", 4);
      ps.python   = ps.newBreed("Python", "snake", 0);
      ps.taipan   = ps.newBreed("Taipan", "snake", 0);

      ps.ralph = ps.newPet("Ralph", ps.labrador, 2014, 02, 11, 4);
      ps.marcy = ps.newPet("Marcy", ps.labrador, 2014, 02, 11, 3);
      ps.boris = ps.newPet("Boris", ps.doberman, 2015, 06, 01, 4);
      ps.guido = ps.newPet("Guido", ps.python  , 1956, 01, 31, 0);
      ps.tracy = ps.newPet("Tracy", ps.taipan  , 2009, 05, 05, 0);

      ps.anneTendsBoris = ps.newTends(ps.anne, ps.boris);
      ps.anneTendsGuido = ps.newTends(ps.anne, ps.guido);
      ps.markTendsRalph = ps.newTends(ps.mark, ps.ralph);
      ps.markTendsMarcy = ps.newTends(ps.mark, ps.marcy);
      ps.markTendsBoris = ps.newTends(ps.mark, ps.boris);

      return ps;
    }

    public PetShop() {
      model = new PetShopModel();
      facade = new PetShopFacade(model);
    }

    public boolean equals(Object obj) {
      if (!(obj instanceof PetShop)) {
        return false;
      }
      return model.equals(((PetShop)obj).model);
    }

    public String toString() {
      return model.toString();
    }

    public Employee newEmployee(String name, int y, int m, int d) {
      LocalDate dob = new LocalDate(y, m, d);
      return facade.getEmployee().insert()
        .setName(name)
        .setDob(dob);
    }
    public Breed newBreed(String name, String species, int numLegs) {
      return facade.getBreed().insert()
        .setBreedName(name)
        .setSpecies(species)
        .setNumLegs(numLegs);
    }

    public Pet newPet(
        String name, Breed breed, int y, int m, int d, int numLegs) {
      LocalDate dob = new LocalDate(y, m, d);
      return facade.getPet().insert()
        .setName(name)
        .setBreed(breed)
        .setDob(dob)
        .setNumLegs(numLegs);
    }

    public EmployeeTendsPet newTends(Employee employee, Pet pet) {
      return facade.getEmployeeTendsPet().build()
        .withEmployee(employee)
        .withPet(pet)
        .insert();
    }
  }

  private void assertModelEquals(PetShop a, PetShop b) {
    assertEquals(a, b);
    if (!a.model.canonicalToString().equals(b.model.canonicalToString())) {
      System.out.println(a.model.canonicalToString());
      System.out.println(b.model.canonicalToString());
    }
    assertEquals(a.model.canonicalToString(), b.model.canonicalToString());
  }

  private void assertModelNotEquals(PetShop a, PetShop b) {
    assertNotEquals(a, b);
    assertNotEquals(a.model.canonicalToString(), b.model.canonicalToString());
  }

  @Test
  public void testEquality() {
    PetShop ps1 = PetShop.base();
    assertModelEquals(ps1, ps1);

    PetShop ps2 = PetShop.base();
    assertModelEquals(ps1, ps2);
    assertModelEquals(ps2, ps1);
  }

  @Test
  public void testSimpleInequality() {
    PetShop ps1 = PetShop.base();
    PetShop ps2 = PetShop.base();

    assertModelEquals(ps1, ps2);

    ps2.ralph.setName("Roberto");
    assertModelNotEquals(ps1, ps2);
  }

  @Test
  public void testEqualityAfterSimpleInequality() {
    PetShop ps1 = PetShop.base();
    PetShop ps2 = PetShop.base();

    assertModelEquals(ps1, ps2);

    ps2.ralph.setName("Roberto");
    ps2.ralph.setName("Ralph");

    assertModelEquals(ps1, ps2);
  }

  @Test
  public void testVaryInsertionOrder() {
    PetShop ps1 = PetShop.base();
    PetShop ps2 = new PetShop();

    ps2.mark = ps2.newEmployee("Mark", 1995, 03, 04);
    ps2.anne = ps2.newEmployee("Anne", 1993, 11, 20);

    ps2.taipan   = ps2.newBreed("Taipan", "snake", 0);
    ps2.labrador = ps2.newBreed("Labrador", "dog", 4);
    ps2.doberman = ps2.newBreed("Doberman Pinscher", "dog", 4);
    ps2.python   = ps2.newBreed("Python", "snake", 0);

    ps2.boris = ps2.newPet("Boris", ps2.doberman, 2015, 06, 01, 4);
    ps2.tracy = ps2.newPet("Tracy", ps2.taipan  , 2009, 05, 05, 0);
    ps2.marcy = ps2.newPet("Marcy", ps2.labrador, 2014, 02, 11, 3);
    ps2.ralph = ps2.newPet("Ralph", ps2.labrador, 2014, 02, 11, 4);
    ps2.guido = ps2.newPet("Guido", ps2.python  , 1956, 01, 31, 0);

    ps2.anneTendsGuido = ps2.newTends(ps2.anne, ps2.guido);
    ps2.markTendsMarcy = ps2.newTends(ps2.mark, ps2.marcy);
    ps2.markTendsBoris = ps2.newTends(ps2.mark, ps2.boris);
    ps2.anneTendsBoris = ps2.newTends(ps2.anne, ps2.boris);
    ps2.markTendsRalph = ps2.newTends(ps2.mark, ps2.ralph);
  }

  @Test
  public void testRemoveThenReinsert() {
    PetShop ps1 = PetShop.base();
    PetShop ps2 = PetShop.base();

    ps2.anneTendsBoris.delete();
    assertModelNotEquals(ps1, ps2);
    assertModelNotEquals(ps1, ps2);
    ps2.anneTendsBoris = ps2.newTends(ps2.anne, ps2.boris);
    assertModelEquals(ps1, ps2);
  }

  @Test
  public void testMinimalRemoveDependentThenReinsert() {
    PetShop ps1 = new PetShop();
    PetShop ps2 = new PetShop();

    ps1.anne = ps1.newEmployee("Anne", 1993, 11, 20);
    ps1.mark = ps1.newEmployee("Mark", 1995, 03, 04);

    ps1.boris = ps1.newPet("Boris", null, 2015, 06, 01, 4);

    ps1.anneTendsBoris = ps1.newTends(ps1.anne, ps1.boris);

    ps2.anne = ps2.newEmployee("Anne", 1993, 11, 20);
    ps2.mark = ps2.newEmployee("Mark", 1995, 03, 04);

    ps2.boris = ps2.newPet("Boris", null, 2015, 06, 01, 4);

    ps2.anneTendsBoris = ps2.newTends(ps2.anne, ps2.boris);

    assertModelEquals(ps1, ps2);
    ps2.anneTendsBoris.delete();
    ps2.anne.delete();
    assertModelNotEquals(ps1, ps2);
    ps2.anne = ps2.newEmployee("Anne", 1993, 11, 20);
    ps2.anneTendsBoris = ps2.newTends(ps2.anne, ps2.boris);
    assertModelEquals(ps1, ps2);
  }

  // @Test
  // public void testRemoveDependentThenReinsert() {
  //   PetShop ps1 = PetShop.base();
  //   PetShop ps2 = PetShop.base();

  //   ps2.anneTendsBoris.delete();
  //   ps2.anneTendsGuido.delete();
  //   ps2.anne.delete();
  //   assertModelNotEquals(ps1, ps2);
  //   ps2.anne = ps2.newEmployee("Anne", 1993, 11, 20);
  //   ps2.anneTendsGuido = ps2.newTends(ps2.anne, ps2.guido);
  //   ps2.anneTendsBoris = ps2.newTends(ps2.anne, ps2.boris);
  //   assertModelEquals(ps1, ps2);
  // }

  @Test
  public void testInsertThenRemove() {
    PetShop ps1 = PetShop.base();
    PetShop ps2 = PetShop.base();

    Breed cockapoo = ps2.newBreed("Cockapoo", "dog", 4);
    assertModelNotEquals(ps1, ps2);
    cockapoo.delete();
    assertModelEquals(ps1, ps2);
  }

}
