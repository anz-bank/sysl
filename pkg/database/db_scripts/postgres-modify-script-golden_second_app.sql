/*TITLE : Petstore Schema*/
/* ---------------------------------------------
Autogenerated script from sysl
--------------------------------------------- */


/*-----------------------Relation Model : RelModelNew-----------------------------------------------*/
CREATE TABLE Customer(
  customerId integer,
  name varchar (50),
  CONSTRAINT CUSTOMER_PK PRIMARY KEY(customerId)
);
