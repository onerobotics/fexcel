/PROG  GET_TASK
/ATTR
COMMENT		= "";
DEFAULT_GROUP	= 1,*,*,*,*;
/APPL
/MN
 : IF (DI[5:ABORT]),JMP LBL[99] ;
 : IF (DI[6:HOME]),JMP LBL[3] ;
 : IF R[3:gripMem]=0,JMP LBL[1] ;
 : IF R[3:gripMem]=1,JMP LBL[2] ;
 : END ;
 :  ;
 : LBL[99] ;
 :   R[1:taskID]=99 ;
 :   END ;
 :  ;
 : LBL[3] ;
 :   R[1:taskID]=3 ;
 :   END ;
 :  ;
 : LBL[1] ;
 :   R[1:taskID]=1 ;
 :   END ;
 :  ;
 : LBL[2] ;
 :   R[1:taskID]=2 ;
 :   END ;
/POS
/END
