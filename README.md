
# Dear Doctor

This is Restfull API build in Golang for booking consultaion slot of various doctors registered. 


## Technologies Used

- Chi router
    : To handle HTTP requests.
- PosgreSQL
    : Used as my database.
- GORM
    : Used to automigrate tables.
- Go-MySQL-Driver
    : To write Sql queries.
- Razor-pay
    : To make payments.
- Docker
    : API and database fully dockerised.
- AWS
    : Hosted in AWS EC2.
- JWT-Token, Pagination, Filltering etc

## Features

- SignUp and SignIn by Admin, Doctor, Users.

- Email verification done for Doctors and Users

- Admin can add new departments.

- Admin can approve Doctors.

- Admin can list all Doctors and Users.

- Admin can fetch all appointments. Can also filter by Day, Doctor.

- Admin can approve Payout requests done by Doctors.

- Doctors can add thier availability.

- Doctors dashboard will show all confirmed appointments.

- Doctors can request for payout.

- Doctors can add thier bank account details for recieving payout.

- Users can make appointments for consultation.

- Users can pay consultation fee using razorpay or COD.







## Environment Variables

To run this project locally, you will need to add the following environment variables to your .env file


To configure database connecton following environment variables are required :-
`DB_HOST`
`DB_DRIVER` 
`DB_USER`
`DB_PASSWORD`
`DB_NAME`
`DB_PORT`

`PORT` : for running locally

`ADMIN_KEY` : To generate token for admin

`DOCTOR_KEY` : To generate token for doctor

`USER_KEY` : To generate token for user

To sent E-mail following environment variables are required :-
`SMTP_PORT`
`SMTP_HOST`
`SMTP_PASSWORD`
`SMTP_USERNAME`

## Screenshots


 ***APIs for Admin***
 
 
<img width="958" alt="Screenshot 2022-09-19 191215" src="https://user-images.githubusercontent.com/49141863/191035381-f87c73f7-24a9-4ae9-b423-8a2e60c69932.png">    

 ***APIs for Doctors***
 
 
<img width="959" alt="Screenshot 2022-09-19 191539" src="https://user-images.githubusercontent.com/49141863/191035436-392a1826-62ac-41e7-a95c-5994bae1e357.png">    

 ***APIs for Users***
 
 
<img width="959" alt="Screenshot 2022-09-19 191703" src="https://user-images.githubusercontent.com/49141863/191035451-380dffd3-20a1-46a5-9592-93e70fcdf602.png">



## Connect with me

Phone : +91 9995 709722

Email : midlaj9995@gmail.com 

LinkedIn: [muhammedali midhilaj](https://www.linkedin.com/in/muhammedali-midhilaj-13834723a)

