/PROG  GET_TASK
/ATTR
COMMENT		= "";
DEFAULT_GROUP	= 1,*,*,*,*;
/APPL
/MN
 : IF (DI{ABORT}),JMP LBL[${TASK_ABORT}] ;
 : IF (DI{HOME}),JMP LBL[${TASK_HOME}] ;
 : IF R{gripMem}=0,JMP LBL[${TASK_PICK}] ;
 : IF R{gripMem}=1,JMP LBL[${TASK_PLACE}] ;
 : END ;
 :  ;
 : LBL[${TASK_ABORT}] ;
 :   R{taskID}=${TASK_ABORT} ;
 :   END ;
 :  ;
 : LBL[${TASK_HOME}] ;
 :   R{taskID}=${TASK_HOME} ;
 :   END ;
 :  ;
 : LBL[${TASK_PICK}] ;
 :   R{taskID}=${TASK_PICK} ;
 :   END ;
 :  ;
 : LBL[${TASK_PLACE}] ;
 :   R{taskID}=${TASK_PLACE} ;
 :   END ;
/POS
/END
