# godb
godb is my attempt to build some kind of sql database

## Based on:
https://cstack.github.io/db_tutorial/

## Todo:
* Btree needs to be implemented kind of 'in file memory'
* Node stores ids (uint32), and those ids represents place in written file where to search for another block

### Binary file look

in vim`:%!xxd`