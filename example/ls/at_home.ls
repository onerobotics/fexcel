/PROG  AT_HOME
/ATTR
COMMENT		= "";
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : UFRAME_NUM=0 ;
 : UTOOL_NUM=0 ;
 : PR[2:LPOS]=LPOS ;
 : IF (PR[2,1]<(-100)),JMP LBL[500] ;
 : IF (PR[2,1]<100),JMP LBL[500] ;
 : IF (PR[2,2]<(-100)),JMP LBL[500] ;
 : IF (PR[2,2]<100),JMP LBL[500] ;
 : IF (PR[2,3]<(-100)),JMP LBL[500] ;
 : IF (PR[2,3]<300),JMP LBL[500] ;
 : R[2:zoneID]=3 ;
 : END ;
 :  ;
 : LBL[500] ;
 :   R[2:zoneID]=0 ;
 :   END ;
/POS
/END
