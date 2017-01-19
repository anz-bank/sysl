import java.lang.reflect.Field;

import io.sysl.test.data.petshop.model.*;

import org.junit.Test;
import static org.junit.Assert.*;

public class TestLazyLoad {

  @Test
  public void testMembersAreNull() {
    try {
      PetShopModel model = new PetShopModel();

      int fieldCount = 0;
      for (Field member : model.getClass().getDeclaredFields()) {
        if (member.getName().startsWith("tbl_")) {
          member.setAccessible(true);
          assertNull(member.get(model));
          fieldCount++;
        }
      }
      assertEquals(fieldCount, 4);
    } catch(Exception e) {
      fail(e.getMessage());
    }
  }

  // TODO(kirkpatg): refactor
  @Test
  public void testMembersAreLazilyConstructed() {
    Field member = null;
    
    try {
      PetShopModel model = new PetShopModel();

      // Breed Table
      Breed.Table breedTable = model.getBreedTable();
      member = model.getClass().getDeclaredField("tbl_Breed");
      member.setAccessible(true);

      assertNotNull(breedTable);
      assertSame(breedTable, member.get(model));

      
      for (Field otherMember : model.getClass().getDeclaredFields()) {
        if (otherMember.getName().startsWith("tbl_")
            && !otherMember.getName().equals("tbl_Breed")) {
          otherMember.setAccessible(true);
          assertNull(otherMember.get(model));
        }
      }

      // Employee Table
      Employee.Table employeeTable = model.getEmployeeTable();
      member = model.getClass().getDeclaredField("tbl_Employee");
      member.setAccessible(true);

      assertNotNull(employeeTable);
      assertSame(employeeTable, member.get(model));

      for (Field otherMember : model.getClass().getDeclaredFields()) {
        if (otherMember.getName().startsWith("tbl_")
            && !otherMember.getName().equals("tbl_Employee")
            && !otherMember.getName().equals("tbl_Breed")) {
          otherMember.setAccessible(true);
          assertNull(otherMember.get(model));
        }
      }

      // EmployeeTendsPet Table
      EmployeeTendsPet.Table etpTable = model.getEmployeeTendsPetTable();
      member = model.getClass().getDeclaredField("tbl_EmployeeTendsPet");
      member.setAccessible(true);

      assertNotNull(etpTable);
      assertSame(etpTable, member.get(model));

      
      for (Field otherMember : model.getClass().getDeclaredFields()) {
        if (otherMember.getName().startsWith("tbl_")
            && !otherMember.getName().equals("tbl_EmployeeTendsPet")
            && !otherMember.getName().equals("tbl_Employee")
            && !otherMember.getName().equals("tbl_Breed")) {
          otherMember.setAccessible(true);
          assertNull(otherMember.get(model));
        }
      }

      // Pet Table
      Pet.Table petTable = model.getPetTable();
      member = model.getClass().getDeclaredField("tbl_Pet");
      member.setAccessible(true);

      assertNotNull(petTable);
      assertSame(petTable, member.get(model));

    } catch(Exception e) {
      fail(e.getMessage());
    }
  }
}
