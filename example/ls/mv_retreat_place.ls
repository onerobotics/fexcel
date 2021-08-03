/PROG  MV_RETREAT_PLACE
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL GET_ZONE ;
 : IF R[2:zoneID]<>2,CALL UNSAFE ;
 :  ;
 : UFRAME_NUM=2 ;
 : UTOOL_NUM=2 ;
 :  ;
 : PR[2:LPOS]=LPOS ;
 : IF (PR[2,3]>250),JMP LBL[1] ;
 :  ;
 : PR[2,3]=250 ;
 : L PR[2:LPOS] 250mm/sec CNT0 ;
 :  ;
 : LBL[1] ;
 : J PR[20:PLACE_PERCH] 100% CNT100 ;
/POS
/END
