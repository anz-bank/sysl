#### Add new commnad reviewdatamodel which can help reviewing generated data model with sysl file produced by commnad import. Generate data model diagrams using the following command
```
sysl reviewdata --root=/Users/guest/data -t Test -o Test.png Test
sysl reviewdatamodel --root=/Users/guest/data -t Test -o Test.png Test.sysl
