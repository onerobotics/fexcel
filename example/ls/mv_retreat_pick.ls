/PROG  MV_RETREAT_PICK
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL GET_ZONE ;
 : IF R[2:zoneID]<>1,CALL UNSAFE ;
 :  ;
 : UFRAME_NUM=1 ;
 : UTOOL_NUM=1 ;
 :  ;
 : PR[2:LPOS]=LPOS ;
 : IF (PR[2,3]>250),JMP LBL[1] ;
 :  ;
 : PR[2,3]=250 ;
 : L PR[2:LPOS] 250mm/sec CNT0 ;
 :  ;
 : LBL[1] ;
 : J PR[10:PICK_PERCH] 100% CNT100 ;
/POS
/END
