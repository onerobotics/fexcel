/PROG  MV_RETREAT_PLACE
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL GET_ZONE ;
 : IF R{zoneID}<>${ZONE_PLACE},CALL UNSAFE ;
 :  ;
 : UFRAME_NUM=${UF_PLACE} ;
 : UTOOL_NUM=${UT_PLACE} ;
 :  ;
 : PR{LPOS}=LPOS ;
 : IF (PR[&PR{LPOS},${Z}]>${PLACE_RETREAT_Z}),JMP LBL[1] ;
 :  ;
 : PR[&PR{LPOS},${Z}]=${PLACE_RETREAT_Z} ;
 : L PR{LPOS} ${RETREAT_SPEED}mm/sec CNT0 ;
 :  ;
 : LBL[1] ;
 : J PR{PLACE_PERCH} ${PERCH_SPEED}% CNT100 ;
/POS
/END
