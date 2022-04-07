<?php

if (isset($_POST['submit'])) {

   //collect form data
   $email = $_POST['email'];

   //if no errors carry on
   if (!isset($error)) {

      // Title of the CSV
      $Content = "Email";

      //set the data of the CSV
      $Content. = "$email";

      //set the file name and create CSV file
      $FileName = "mailing_list1.csv"
      header('Content-Type: application/csv');
      header('Content-Disposition: attachment; filename="'.$FileName.
         '"');
      echo $Content;
      exit();
   }
}

//if their are errors display them
if (isset($error)) {
   foreach($error as $error) {
      echo '$error';
   }
}


?>
