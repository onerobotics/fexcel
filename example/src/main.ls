/PROG  MAIN
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL INIT ;
 : CALL SV_HOME ;
 :  ;
 : LBL[1] ;
 :   CALL GET_TASK ;
 :   SELECT R{taskID}=${TASK_PICK},CALL SV_PICK ;
 :          =${TASK_PLACE},CALL SV_PLACE ;
 :          =${TASK_HOME},CALL SV_HOME ;
 :          =${TASK_ABORT},CALL SV_ABORT ;
 :   JMP LBL[1] ;
/POS
/END
