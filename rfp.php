<?php

class rfp extends Customer_Controller{

	public function __construct(){
        parent::__construct();
        
		$this->data['showAdminMenu'] = 1;
		$this->data['sideBar'] = 'rooms';
        $this->load->model('employee_m');
    }

	public function index(){

		 $this->data['subview'] = 'customer/rfp/index';
		
		incView($this->data['subview'], $this->data);
	}
	public function hotel(){

		 $this->data['subview'] = 'customer/rfp/index';
		
		incView($this->data['subview'], $this->data);
	}
	public function curl($url,$data=null,$headers=null){
		$ch = curl_init($url);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
	
		if(!empty($data)){
			curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
		}
	
		if (!empty($headers)) {
			curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
		}
	
		$response = curl_exec($ch);
	
		if (curl_error($ch)) {
			trigger_error('Curl Error:' . curl_error($ch));
		}
	
		curl_close($ch);
		return $response;
	}

	public function getQues(){
		if(isset($_POST)){
			$id= $this->session->userdata['id'];
			$ques = $_POST['questionCategoryParent'];
			$res = $this->curl('http://localhost:9000/ansHotel/edit','questionCategoryParent=1&travelAgencyMasterId='.$id);
			echo $res;
		}
	}
	public function sendQues(){
		if(isset($_POST)){
			$_POST['travelAgencyMasterId'] =  $this->session->userdata['id'];
			$_POST['clientTypeMasterId'] = '1';
			$data = json_encode($_POST);
			$res = $this->curl('http://localhost:9000/ansHotel',$data);echo $res;
			//echo $data;
		}
	}
	


}