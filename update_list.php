<?php

if (isset($_POST['email'])) {

    $fA = $_POST["email"];
    //$fB = $_POST["name"];
    $keys = array($fA);  //, $fB);
    $csv_line = $keys;
    foreach( $keys as $key ){
        array_push($csv_line,'' . $_GET[$key]);
    }
    $csv_line = implode(',',$csv_line);
    $fname = 'mailing_list.csv';
    if(!file_exists($fname)){$csv_line = $csv_line."\r\n" ;}
    $fp = fopen($fname,'a');
    //print_r(error_get_last());
    //$fcontent = $csv_line;
    fwrite($fp,$csv_line);
    fclose($fp);
    //echo $csv_line;
    echo "Thank you for subscribing to the mailing list!";
} else {
    echo "Empty field submitted.";
}

?>
