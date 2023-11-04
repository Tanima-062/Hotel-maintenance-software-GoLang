package utils

const (
	// WholesalerIDParent 親アカウントに便宜的につけているホールセラーID
	WholesalerIDParent = 0
	// WholesalerIDTl TLのホールセラーID
	WholesalerIDTl = 3
	// WholesalerIDTema TemaのホールセラーID
	WholesalerIDTema = 4
	// WholesalerIDNeppan NeppanのホールセラーID
	WholesalerIDNeppan = 6
	// WholesalerIDDirect 直仕入れのホールセラーID
	WholesalerIDDirect = 7
	// WholesalerIDRaku2 らく通のホールセラーID
	WholesalerIDRaku2 = 8

	// ImgBasePath gcsアップロードの際のbase url
	ImgBasePath = "hotel/property-images"

	// ChildRateTypeA ChildAに該当するchild_rate_type
	ChildRateTypeA = 1
	// ChildRateTypeB ChildBに該当するchild_rate_type
	ChildRateTypeB = 2
	// ChildRateTypeC ChildCに該当するchild_rate_type
	ChildRateTypeC = 3
	// ChildRateTypeD ChildDに該当するchild_rate_type
	ChildRateTypeD = 4
	// ChildRateTypeE ChildEに該当するchild_rate_type
	ChildRateTypeE = 5
	// ChildRateTypeF ChildFに該当するchild_rate_type
	ChildRateTypeF = 6

	// ReserveStatusReserved 予約済み
	ReserveStatusReserved = 1
	// ReserveStatusCancel 予約キャンセル
	ReserveStatusCancel = 2
	// ReserveStatusNoShow 予約NoShow
	ReserveStatusNoShow = 3
	// ReserveStatusStaying 宿泊中
	ReserveStatusStaying = 4
	// ReserveStatusStayed 旅程終了
	ReserveStatusStayed = 5

	LogServiceRoom      = "ROOM"
	LogServiceStock     = "STOCK"
	LogServicePlan      = "PLAN"
	LogServicePrice     = "PRICE"
	LogTypeMaster       = "Master"
	LogTypeDifferential = "Differential"

	// TlApiHeader XML API通信時のヘッダーテンプレート
	TlApiHeader = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:head="http://www.seanuts.co.jp/ota/header" xmlns:ns="http://www.opentravel.org/OTA/2003/05">
	<soapenv:Header>
		<head:Interface>
			<head:PayloadInfo>
				<head:CommDescriptor>
					<head:Authentication>
						<head:Username>%s</head:Username>
						<head:Password>%s</head:Password>
					</head:Authentication>
				</head:CommDescriptor>
			</head:PayloadInfo>
		</head:Interface>
	</soapenv:Header>
	<soapenv:Body>
	%s
	</soapenv:Body>
	</soapenv:Envelope>`

	// TlAPIReadRequstXML API通信時のボディテンプレート
	TlAPIReadRequst = `<ns:OTA_ReadRQ Version="1.0" PrimaryLangID="jpn">
	<ns:POS>
		<ns:Source>
			<ns:RequestorID Type="5">
				<ns:CompanyName>%s</ns:CompanyName>
			</ns:RequestorID>
		</ns:Source>
	</ns:POS>
	<ns:ReadRequests>
		<ns:ReadRequest>
			<ns:UniqueID Type="%s" ID="%s" ID_Context="%s"/>
			<ns:Verification>
				<ns:TPA_Extensions>
					<ns:BasicPropertyInfo HotelCode="%s" />
				</ns:TPA_Extensions>
			</ns:Verification>
		</ns:ReadRequest>
	</ns:ReadRequests>
</ns:OTA_ReadRQ>`
)
