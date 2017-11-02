package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var jsonString []byte

func AddRow(values []string, columns []string, table string) int64 {
	db = GetDB()
	// insert
	column := strings.Join(columns, ",")
	value := strings.Join(values, "','")
	stmt, err := db.Exec("INSERT into " + table + "(" + column + ") VALUES ('" + value + "')")
	checkErr(err)

	rowid, err := stmt.LastInsertId()
	checkErr(err)

	fmt.Println("Value Inserted into ", table)
	//fmt.Println("RowId ", rowid)

	return rowid

}

func EditRow(formVal []string, columnVal []string, formCondVal []string, columnCondVal []string, table string) int64 {
	db = GetDB()
	//column := strings.Join(columns, ",")
	//value := strings.Join(values, "','")
	setStr := "SET"
	//fmt.Println(len(columnVal))
	i := 0
	for i = 0; i < len(columnVal); i++ {
		setStr = setStr + " , " + columnVal[i] + " = '" + formVal[i] + "'"

	}
	setStr = strings.Replace(setStr, "SET ,", " SET", -2)
	//fmt.Println(setStr)
	conditionStr := createCondStr(formCondVal, columnCondVal)
	//fmt.Println("UPDATE " + table + setStr + conditionStr)
	stmt, err := db.Exec("UPDATE " + table + setStr + conditionStr)

	checkErr(err)

	rowcnt, err := stmt.RowsAffected()
	checkErr(err)
	//fmt.Println("Row(s) Updated ", rowcnt)

	return rowcnt
}

func QueryRow(values []string, columns []string, table string) int {
	db = GetDB()
	conditionStr := createCondStr(values, columns)

	// select
	var cnt int
	//fmt.Println(">>>>>>>>>>>>", conditionStr)
	_ = db.QueryRow("select count(*) from " + table + conditionStr).Scan(&cnt)
	//fmt.Println(">>>>", cnt)
	//fmt.Println(cnt)
	return cnt
}

func QueryRowEdit(values []string, columns []string, reqid, table string) int {
	db = GetDB()
	conditionStr := createCondStr(values, columns)

	// select
	var cnt int
	//fmt.Println(">>>>>>>>>>>>", conditionStr)
	_ = db.QueryRow("select count(*) from " + table + conditionStr + "AND designationMasterId != " + reqid).Scan(&cnt)
	//fmt.Println(">>>>", cnt)
	//fmt.Println(cnt)
	return cnt
}

func ListAllCity() string {
	db = GetDB()
	var cities []LabVal
	var city LabVal
	//var label string
	stmt, err := db.Query("SELECT cityMasterId as value, cityName as label from citymaster")
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&city.Value, &city.Label)
		checkErr(err)

		city = LabVal{
			Label: city.Label,
			Value: city.Value,
		}
		cities = append(cities, city)
	}
	b, err := json.Marshal(cities)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}

func ListCityCat(companyId string) string {
	db = GetDB()

	stmt, err := db.Query("SELECT CityCatName as label, CityCatID as value from city_category where CompanyID = '" + companyId + "'")
	checkErr(err)
	return ParseRow(stmt)
}

func ListCity(citycatId string) string {
	db = GetDB()
	var cities []LabVal
	var labval LabVal
	var cityCatName string
	//var label string
	stmt, err := db.Query("SELECT cmas.cityName, cmap.CityID from city_mapping as cmap JOIN citymaster as cmas ON cmas.cityMasterId = cmap.CityID where cmap.CityCatID = '" + citycatId + "'")
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&labval.Label, &labval.Value)
		checkErr(err)

		labval = LabVal{
			Label: labval.Label,
			Value: labval.Value,
		}
		cities = append(cities, labval)
	}
	//fmt.Println(cities)
	err = db.QueryRow("SELECT CityCatName from city_category where CityCatID = '" + citycatId + "'").Scan(&cityCatName)
	checkErr(err)

	labval = LabVal{
		Label: cityCatName,
		Value: citycatId,
	}
	cityCategoryMap := CityCategoryMap{
		CompanyID: "",
		CityCat:   labval,
		Cities:    cities,
	}
	b, err := json.Marshal(cityCategoryMap)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}
func ListBundleRequirements(tablename, companyId string) string {
	db = GetDB()
	//var labval LabVal
	var benefittype LabVal
	var policybundle PolicyBundle
	var mybundles []PolicyBundle
	var cityCatAndAllowance CityCatAndAllowance

	stmt, err := db.Query("SELECT BenefitTypeID, BenefitTypeName from benefit_type_master")
	checkErr(err)

	for stmt.Next() {
		err := stmt.Scan(&benefittype.Value, &benefittype.Label)
		checkErr(err)
		benefittype = LabVal{
			Label: benefittype.Label,
			Value: benefittype.Value,
		}
		stmt1, err := db.Query("SELECT BenefitID, BenefitName FROM benefit_master WHERE BenefitTypeID='" + benefittype.Value + "'")

		var benefit LabVal
		var cityCatAndAllowances []CityCatAndAllowance
		//var benefit LabVal
		var benefits []LabVal
		for stmt1.Next() {
			err := stmt1.Scan(&benefit.Value, &benefit.Label)
			checkErr(err)
			benefit = LabVal{
				Label: benefit.Label,
				Value: benefit.Value,
			}
			benefits = append(benefits, benefit)
		}
		stmt2, err := db.Query("SELECT CityCatName, CityCatID from city_category where CompanyID = '" + companyId + "'")
		for stmt2.Next() {
			err := stmt2.Scan(&cityCatAndAllowance.Label, &cityCatAndAllowance.Value)
			checkErr(err)
			cityCatAndAllowance = CityCatAndAllowance{
				Label:      cityCatAndAllowance.Label,
				Value:      cityCatAndAllowance.Value,
				LimitSpent: "false",
				Min:        "",
				Max:        "",
				Flex:       "",
				FlexAmt:    "",
				StarCat:    "",
			}
			cityCatAndAllowances = append(cityCatAndAllowances, cityCatAndAllowance)
			//benefits = append(benefits, labval)
		}
		policybundle = PolicyBundle{
			BenefitTypeID:        benefittype,
			Priority:             policybundle.Priority,
			Benefits:             benefits,
			CityCatAndAllowances: cityCatAndAllowances,
		}
		mybundles = append(mybundles, policybundle)
	}
	//var pb PB

	var pb = PB{
		BenefitBundleID: "",
		BundleName:      "",
		BundleCode:      "",
		CompanyID:       companyId,
		MethType:        "",
		PolicyBundles:   mybundles,
	}
	b, err := json.Marshal(pb)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}
func ListBundleDetail(tablename string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	//conditionStr := createCondStr(formCondVal, columnCondVal)
	var pb PB
	var policybundle PolicyBundle
	var labval LabVal
	var benefittype LabVal
	var cityCatAndAllowance CityCatAndAllowance
	var mybundles []PolicyBundle
	var benefitBundleTypeMappingId string

	conditionStr := createCondStr(formCondVal, columnCondVal)

	stmt, err := db.Query("select BenefitBundleID, BenefitBundleName, BenefitBundleCode, CompanyID from policy_benefit_bundle " + conditionStr)
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&pb.BenefitBundleID, &pb.BundleName, &pb.BundleCode, &pb.CompanyID)
		checkErr(err)
	}
	stmt1, err := db.Query("SELECT bbtm.BenefitBundleTypeMappingID, bbtm.BenefitTypeID, btm.BenefitTypeName, bbtm.Priority from benefit_bundle_type_mapping as bbtm JOIN benefit_type_master as btm ON bbtm.BenefitTypeID = btm.BenefitTypeID WHERE bbtm.BenefitBundleID = '" + pb.BenefitBundleID + "'")
	checkErr(err)

	for stmt1.Next() {
		err := stmt1.Scan(&benefitBundleTypeMappingId, &labval.Value, &labval.Label, &policybundle.Priority)
		checkErr(err)

		benefittype = LabVal{
			Label: labval.Label,
			Value: labval.Value,
		}
		stmt2, err := db.Query("SELECT btbm.BenefitID, bm.BenefitName FROM bundle_type_benefit_mapping as btbm JOIN benefit_master as bm ON btbm.BenefitID = bm.BenefitID  WHERE BenefitBundleTypeMappingID='" + benefitBundleTypeMappingId + "'")

		var cityCatAndAllowances []CityCatAndAllowance
		var benefits []LabVal
		for stmt2.Next() {
			err := stmt2.Scan(&labval.Value, &labval.Label)
			checkErr(err)
			labval = LabVal{
				Label: labval.Label,
				Value: labval.Value,
			}
			benefits = append(benefits, labval)
		}
		stmt3, err := db.Query("SELECT ct.CityCatName, btam.CityCatID, btam.LimitSpend, btam.MaxAmount, btam.MinAmount, btam.Flexibility, btam.FlexAmount, btam.StarCat FROM benefit_type_allowance_mapping as btam JOIN city_category as ct ON btam.CityCatID = ct.CityCatID WHERE BenefitBundleTypeMappingID = '" + benefitBundleTypeMappingId + "'")
		for stmt3.Next() {
			err := stmt3.Scan(&cityCatAndAllowance.Label, &cityCatAndAllowance.Value, &cityCatAndAllowance.LimitSpent, &cityCatAndAllowance.Max, &cityCatAndAllowance.Min, &cityCatAndAllowance.Flex, &cityCatAndAllowance.FlexAmt, &cityCatAndAllowance.StarCat)
			checkErr(err)
			cityCatAndAllowance = CityCatAndAllowance{
				Label:      cityCatAndAllowance.Label,
				Value:      cityCatAndAllowance.Value,
				LimitSpent: cityCatAndAllowance.LimitSpent,
				Min:        cityCatAndAllowance.Min,
				Max:        cityCatAndAllowance.Max,
				Flex:       cityCatAndAllowance.Flex,
				FlexAmt:    cityCatAndAllowance.FlexAmt,
				StarCat:    cityCatAndAllowance.StarCat,
			}
			cityCatAndAllowances = append(cityCatAndAllowances, cityCatAndAllowance)
			//benefits = append(benefits, labval)
		}

		policybundle = PolicyBundle{
			BenefitTypeID:        benefittype,
			Priority:             policybundle.Priority,
			Benefits:             benefits,
			CityCatAndAllowances: cityCatAndAllowances,
		}
		mybundles = append(mybundles, policybundle)

	}

	pb = PB{
		BenefitBundleID: pb.BenefitBundleID,
		BundleName:      pb.BundleName,
		BundleCode:      pb.BundleCode,
		CompanyID:       pb.CompanyID,
		MethType:        "",
		PolicyBundles:   mybundles,
	}
	b, err := json.Marshal(pb)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}

func ListRow(tablename string, companyId string) string {
	db = GetDB()

	stmt, err := db.Query("select * from " + tablename + " where CompanyID = " + companyId)
	checkErr(err)
	return ParseRow(stmt)
}

func ParseRow(stmt *sql.Rows) string {

	columnNames, err := stmt.Columns()
	checkErr(err)
	rc := NewMapStringScan(columnNames)
	var slice = "["
	for stmt.Next() {
		err := rc.Update(stmt)
		checkErr(err)
		cv := rc.Get()
		//log.Printf("%#v\n\n >>>>>>", cv[columnNames[0]])
		jsonString, _ = json.Marshal(cv)

		slice += string(jsonString) + ","
		//fmt.Println(string(slice))

	}
	slice += "]"
	r := strings.NewReplacer(",]", "]")
	slice = r.Replace(slice)
	return slice
}

func CityMapCheck(cityId, companyId string) string {
	db = GetDB()
	// insert
	var citycat = ""
	_ = db.QueryRow("select CityCatID from city_mapping where cityID = '" + cityId + "' AND CityCatID IN (SELECT CityCatID from city_category where companyID = '" + companyId + "') ").Scan(&citycat)

	_ = db.QueryRow("select CityCatName from city_category where CityCatID = '" + citycat + "'").Scan(&citycat)

	fmt.Println("/////", citycat)
	//rowcnt, err := stmt.RowsAffected()
	//checkErr(err)
	//fmt.Println(rowcnt)
	return citycat
}

// func getBenefitBundleID(benefitBundleCode string, companyId string) string {
// 	db = GetDB()
// 	var benefitBundleId string
// 	_ = db.QueryRow("select BenefitBundleID from policy_benefit_bundle where BenefitBundleCode = '" + benefitBundleCode + "' and CompanyID = '" + companyId + "'").Scan(&benefitBundleId)
// 	return benefitBundleId
// }
func ListBundles(table, companyId string) string {
	db = GetDB()
	var bundle LabVal
	var bundlelist []LabVal
	stmt, err := db.Query("select BenefitBundleName as label, BenefitBundleID as value from " + table + " where CompanyID = '" + companyId + "'")
	checkErr(err)

	for stmt.Next() {
		err := stmt.Scan(&bundle.Label, &bundle.Value)
		checkErr(err)
		bundle = LabVal{
			Label: bundle.Label,
			Value: bundle.Value,
		}
		bundlelist = append(bundlelist, bundle)
	}
	b, err := json.Marshal(bundlelist)
	return string(b)
}

func ListDepartmentDetail(tablename string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	var departmentVar Department
	conditionStr := createCondStr(formCondVal, columnCondVal)

	stmt, err := db.Query("select departmentMasterId, departmentName, departmentCode, travelAgencyMasterId from " + tablename + conditionStr)
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&departmentVar.DepartmentID, &departmentVar.DepartmentName, &departmentVar.DepartmentCode, &departmentVar.TravelAgencyMasterID)
		checkErr(err)
		departmentVar = Department{
			DepartmentID:         departmentVar.DepartmentID,
			DepartmentName:       departmentVar.DepartmentName,
			DepartmentCode:       departmentVar.DepartmentCode,
			TravelAgencyMasterID: departmentVar.TravelAgencyMasterID,
		}
	}
	b, err := json.Marshal(departmentVar)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}

func ListDepartments(table, companyId string) string {
	db = GetDB()
	var bundle LabVal
	var bundlelist []LabVal
	fmt.Println(companyId + "  " + table)
	stmt, err := db.Query("select departmentName as label, departmentMasterId as value from " + table + " where travelAgencyMasterId = '" + companyId + "'")
	checkErr(err)

	for stmt.Next() {
		err := stmt.Scan(&bundle.Label, &bundle.Value)
		checkErr(err)
		bundle = LabVal{
			Label: bundle.Label,
			Value: bundle.Value,
		}
		bundlelist = append(bundlelist, bundle)
	}
	b, err := json.Marshal(bundlelist)
	return string(b)
}

func ListDesignationsByDep(table, companyId string, formCondVal, columnCondVal []string) string {
	db = GetDB()
	var bundle LabVal
	var bundlelist []LabVal
	fmt.Println(companyId + "  " + table)

	conditionStr := createCondStr(formCondVal, columnCondVal)

	stmt, err := db.Query("select designationName as label, designationMasterId as value from " + table + conditionStr + " AND travelAgencyMasterId = '" + companyId + "'")
	checkErr(err)

	for stmt.Next() {
		err := stmt.Scan(&bundle.Label, &bundle.Value)
		checkErr(err)
		bundle = LabVal{
			Label: bundle.Label,
			Value: bundle.Value,
		}
		bundlelist = append(bundlelist, bundle)
	}
	b, err := json.Marshal(bundlelist)
	return string(b)
}

func ListDesignaionDetail(tablename string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	var designationVar Designation
	conditionStr := createCondStr(formCondVal, columnCondVal)

	stmt, err := db.Query("select designationMasterId, designationName, designationCode, hierarchyId, travelAgencyMasterId, benefitBundleId, department from " + tablename + conditionStr)
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&designationVar.DesignationID, &designationVar.DesignationName, &designationVar.DesignationCode, &designationVar.HierarchyID, &designationVar.TravelAgencyMasterID, &designationVar.BenefitBundleID, &designationVar.Department)
		checkErr(err)
		designationVar = Designation{
			DesignationID:        designationVar.DesignationID,
			DesignationName:      designationVar.DesignationName,
			DesignationCode:      designationVar.DesignationCode,
			HierarchyID:          designationVar.HierarchyID,
			TravelAgencyMasterID: designationVar.TravelAgencyMasterID,
			BenefitBundleID:      designationVar.BenefitBundleID,
			Department:           designationVar.Department,
		}
	}
	b, err := json.Marshal(designationVar)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}

func ListDesignations(table, companyId string) string {
	db = GetDB()
	var bundle LabVal
	var bundlelist []LabVal
	//fmt.Println(companyId + "  " + table)
	stmt, err := db.Query("select designationName as label, designationMasterId as value from " + table + " where travelAgencyMasterId = '" + companyId + "'")
	checkErr(err)

	for stmt.Next() {
		err := stmt.Scan(&bundle.Label, &bundle.Value)
		checkErr(err)
		bundle = LabVal{
			Label: bundle.Label,
			Value: bundle.Value,
		}
		bundlelist = append(bundlelist, bundle)
	}
	b, err := json.Marshal(bundlelist)
	return string(b)
}

func GetEmployeedetails(table, travelAgencyUsersId string) string {
	db = GetDB()
	var editEmployeeVar EditEmploye
	var benefitbundle LabVal
	var designation LabVal
	var department LabVal
	stmt, err := db.Query("select travelAgencyUsersId, travelAgencyNameTemp, virtualName, email, personalEmail, phone, mobile, designation, designationId, hierarchyId, benefitBundleId from " + table + " where travelAgencyUsersId = '" + travelAgencyUsersId + "'")
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&editEmployeeVar.TravelAgencyUserId, &editEmployeeVar.CompanyName, &editEmployeeVar.VirtualName, &editEmployeeVar.Email, &editEmployeeVar.PersonalEmail, &editEmployeeVar.Phone, &editEmployeeVar.Mobile, &editEmployeeVar.Designation.Label, &editEmployeeVar.Designation.Value, &editEmployeeVar.HierarchyId, &editEmployeeVar.BenefitBundle.Value)
		checkErr(err)
	}

	err = db.QueryRow("select department from designationmaster where designationMasterId = '" + editEmployeeVar.Designation.Value + "'").Scan(&editEmployeeVar.Department.Value)
	checkErr(err)
	err = db.QueryRow("select departmentName from department where departmentMasterId = '" + editEmployeeVar.Department.Value + "'").Scan(&editEmployeeVar.Department.Label)
	checkErr(err)

	department = LabVal{
		Label: editEmployeeVar.Department.Label,
		Value: editEmployeeVar.Department.Value,
	}

	err = db.QueryRow("select BenefitBundleName from policy_benefit_bundle where BenefitBundleID = '" + editEmployeeVar.BenefitBundle.Value + "'").Scan(&editEmployeeVar.BenefitBundle.Label)
	checkErr(err)

	benefitbundle = LabVal{
		Label: editEmployeeVar.BenefitBundle.Label,
		Value: editEmployeeVar.BenefitBundle.Value,
	}

	designation = LabVal{
		Label: editEmployeeVar.Designation.Label,
		Value: editEmployeeVar.Designation.Value,
	}
	editEmployeeVar = EditEmploye{
		TravelAgencyUserId: travelAgencyUsersId,
		CompanyName:        editEmployeeVar.CompanyName,
		VirtualName:        editEmployeeVar.VirtualName,
		Email:              editEmployeeVar.Email,
		PersonalEmail:      editEmployeeVar.PersonalEmail,
		HierarchyId:        editEmployeeVar.HierarchyId,
		Phone:              editEmployeeVar.Phone,
		Mobile:             editEmployeeVar.Mobile,
		Department:         department,
		Designation:        designation,
		BenefitBundle:      benefitbundle,
	}

	b, err := json.Marshal(editEmployeeVar)
	return string(b)
}

func SeachAllEmp(table, travelAgencyMasterId, travelAgencyUserName string) string {
	db = GetDB()
	var allEmps []AllEmp
	stmt, err := db.Query("select travelAgencyUsersId, virtualName, email, personalEmail, phone, hierarchyId, designationId, benefitBundleId from " + table + " where travelAgencyMasterId = '" + travelAgencyMasterId + "' and virtualName like '%" + travelAgencyUserName + "%'")
	checkErr(err)

	var benefitbundle LabVal
	var designation LabVal
	var department LabVal
	for stmt.Next() {
		var allEmp AllEmp

		err := stmt.Scan(&allEmp.TravelAgencyUserId, &allEmp.VirtualName, &allEmp.Email, &allEmp.PersonalEmail, &allEmp.Phone, &allEmp.HierarchyId, &allEmp.Designation.Value, &allEmp.BenefitBundle.Value)
		checkErr(err)
		err = db.QueryRow("select BenefitBundleName from policy_benefit_bundle where BenefitBundleID = '" + allEmp.BenefitBundle.Value + "'").Scan(&allEmp.BenefitBundle.Label)
		checkErr(err)

		err = db.QueryRow("select department, designationName from designationmaster where designationMasterId = '"+allEmp.Designation.Value+"'").Scan(&allEmp.Department.Value, &allEmp.Designation.Label)
		checkErr(err)

		err = db.QueryRow("select departmentName from department where departmentMasterId = '" + allEmp.Department.Value + "'").Scan(&allEmp.Department.Label)
		checkErr(err)

		fmt.Println(allEmp.BenefitBundle.Label, "   ", allEmp.VirtualName)
		benefitbundle = LabVal{
			Label: allEmp.BenefitBundle.Label,
			Value: allEmp.BenefitBundle.Value,
		}
		designation = LabVal{
			Label: allEmp.Designation.Label,
			Value: allEmp.Designation.Value,
		}
		department = LabVal{
			Label: allEmp.Department.Label,
			Value: allEmp.Department.Value,
		}
		allEmp = AllEmp{
			TravelAgencyUserId: allEmp.TravelAgencyUserId,
			VirtualName:        allEmp.VirtualName,
			Email:              allEmp.Email,
			PersonalEmail:      allEmp.PersonalEmail,
			Phone:              allEmp.Phone,
			HierarchyId:        allEmp.HierarchyId,
			Department:         department,
			Designation:        designation,
			BenefitBundle:      benefitbundle,
		}
		allEmps = append(allEmps, allEmp)
	}
	b, err := json.Marshal(allEmps)
	return string(b)
}

func ListAllEmployees(table, companyID string, FormCondVal, ColumnCondVal []string) string {
	db = GetDB()
	var allEmps []AllEmp
	conditionStr := createCondStr(FormCondVal, ColumnCondVal)
	stmt, err := db.Query("select travelAgencyUsersId, virtualName, email, personalEmail, phone, hierarchyId, designationId, benefitBundleId from " + table + conditionStr + " and designationId IS NOT NULL and status = '1'")
	checkErr(err)
	for stmt.Next() {
		var allEmp AllEmp

		err := stmt.Scan(&allEmp.TravelAgencyUserId, &allEmp.VirtualName, &allEmp.Email, &allEmp.PersonalEmail, &allEmp.Phone, &allEmp.HierarchyId, &allEmp.Designation.Value, &allEmp.BenefitBundle.Value)
		checkErr(err)

		err = db.QueryRow("select BenefitBundleName from policy_benefit_bundle where BenefitBundleID = '" + allEmp.BenefitBundle.Value + "'").Scan(&allEmp.BenefitBundle.Label)
		checkErr(err)

		err = db.QueryRow("select department, designationName from designationmaster where designationMasterId = '"+allEmp.Designation.Value+"'").Scan(&allEmp.Department.Value, &allEmp.Designation.Label)
		checkErr(err)

		err = db.QueryRow("select departmentName from department where departmentMasterId = '" + allEmp.Department.Value + "'").Scan(&allEmp.Department.Label)
		checkErr(err)

		var benefitbundle LabVal
		var designation LabVal
		var department LabVal

		benefitbundle = LabVal{
			Label: allEmp.BenefitBundle.Label,
			Value: allEmp.BenefitBundle.Value,
		}
		designation = LabVal{
			Label: allEmp.Designation.Label,
			Value: allEmp.Designation.Value,
		}
		department = LabVal{
			Label: allEmp.Department.Label,
			Value: allEmp.Department.Value,
		}

		allEmp = AllEmp{
			TravelAgencyUserId: allEmp.TravelAgencyUserId,
			VirtualName:        allEmp.VirtualName,
			Email:              allEmp.Email,
			PersonalEmail:      allEmp.PersonalEmail,
			Phone:              allEmp.Phone,
			HierarchyId:        allEmp.HierarchyId,
			Department:         department,
			Designation:        designation,
			BenefitBundle:      benefitbundle,
		}
		allEmps = append(allEmps, allEmp)
	}
	b, err := json.Marshal(allEmps)
	return string(b)

}
func ListUnassignedEmployees(table, companyID string, FormCondVal, ColumnCondVal []string) string {
	db = GetDB()
	var empDetail EmpDetail
	var empDetails []EmpDetail

	conditionStr := createCondStr(FormCondVal, ColumnCondVal)

	//var Id string
	stmt, err := db.Query("select travelAgencyUsersId, virtualName, email, personalEmail, phone, designation from " + table + conditionStr + " and designationId IS NULL and status = '1'")
	//return Id
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&empDetail.TravelAgencyUserId, &empDetail.VirtualName, &empDetail.Email, &empDetail.PersonalEmail, &empDetail.Phone, &empDetail.Designation)
		checkErr(err)
		empDetail = EmpDetail{
			TravelAgencyUserId: empDetail.TravelAgencyUserId,
			VirtualName:        empDetail.VirtualName,
			Email:              empDetail.Email,
			PersonalEmail:      empDetail.PersonalEmail,
			Phone:              empDetail.Phone,
			Designation:        empDetail.Designation,
		}
		empDetails = append(empDetails, empDetail)
	}
	b, err := json.Marshal(empDetails)
	return string(b)

}

func getId(table, requestId string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	conditionStr := createCondStr(formCondVal, columnCondVal)

	var Id string
	_ = db.QueryRow("select " + requestId + " from " + table + conditionStr).Scan(&Id)
	return Id
}

func checkDependency(cityCatId string) int {
	db = GetDB()
	var count int
	_ = db.QueryRow("SELECT COUNT(1) FROM `benefit_type_allowance_mapping` WHERE CityCatID = " + cityCatId).Scan(&count)
	return count
}

// func getBenefitBundleTypeMappingID(benefitBundleType string, benefitBundleID string) string {
// 	db = GetDB()
// 	var benefitBundleTypeMappingId string
// 	_ = db.QueryRow("select BenefitBundleTypeMappingID from benefit_bundle_type_mapping where BenefitTypeID = '" + benefitBundleType + "' and BenefitBundleID = '" + benefitBundleID + "'").Scan(&benefitBundleTypeMappingId)
// 	return benefitBundleTypeMappingId
// }

// func GetDependency(tablename string) string {
// 	fmt.Println("I am inside get dependency")
// 	stmt, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = 'company_policy' AND REFERENCED_TABLE_NAME = '" + tablename + "'")
// 	checkErr(err)
// 	fmt.Println("ran query")

// 	var tables string
// 	fmt.Println("all set! go!")

// 	for stmt.Next() {
// 		var name string

// 		if err := stmt.Scan(&name); err != nil {
// 			log.Fatal(err)
// 		}
// 		tables = tables + "," + name
// 	}
// 	fmt.Printf("......", tables)
// 	return tables
// }

func DeleteById(tablename string, formCondVal []string, columnCondVal []string) string {

	db = GetDB()
	//_, err := db.Exec("delete from " + tablename + createCondStr(formCondVal, columnCondVal))
	//checkErr(err)
	fmt.Println(createCondStr(formCondVal, columnCondVal))
	stmt, err := db.Exec("delete from " + tablename + createCondStr(formCondVal, columnCondVal))
	checkErr(err)

	affect, err := stmt.RowsAffected()
	checkErr(err)
	if affect > 0 {
		return string(affect) + " records deleted"
	} else {
		return "Record not exists"
	}

	//	return string(affect)

	//fmt.Println(affect)

}

func createCondStr(formCondVal []string, columnCondVal []string) string {
	conditionStr := " WHERE"
	i := 0
	for i = 0; i < len(columnCondVal); i++ {
		conditionStr = conditionStr + " and " + columnCondVal[i] + " = '" + formCondVal[i] + "'"
	}
	conditionStr = strings.Replace(conditionStr, "WHERE and", " WHERE", -2)
	return conditionStr
}

func GetDB() *sql.DB {
	var err error
	if db == nil {
		//db, err = sql.Open("mysql", "root:@/company_policy?parseTime=true&charset=utf8")
		db, err = sql.Open("mysql", "sriram:sriram123@tcp(127.0.0.1:3306)/hotnix_dev?charset=utf8")

		checkErr(err)
	}

	return db
}
