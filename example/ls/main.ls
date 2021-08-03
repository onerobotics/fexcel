/PROG  MAIN
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL INIT ;
 : CALL SV_HOME ;
 :  ;
 : LBL[1] ;
 :   CALL GET_TASK ;
 :   SELECT R[1:taskID]=1,CALL SV_PICK ;
 :          =2,CALL SV_PLACE ;
 :          =3,CALL SV_HOME ;
 :          =99,CALL SV_ABORT ;
 :   JMP LBL[1] ;
/POS
/END
