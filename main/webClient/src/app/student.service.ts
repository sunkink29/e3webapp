import { Injectable } from '@angular/core';
import { Observable ,  of ,  Observer ,  Subscription } from 'rxjs';
import { HttpClient } from '@angular/common/http';

import { Student, Block } from './Interfaces';
import { TeacherService } from './teacher.service';
import { DataContainer, getHttpCacheObervable } from './HttpCacheObervable'

@Injectable()
export class StudentService {

  private makeStudentObservable(index: number, url: string,
    dataContainer: DataContainer<Student[][], Student[]>, fakeData: Student[][]) {
    return getHttpCacheObervable<Student[][], Student[], Student[][], Student[][]>(
      this.http, dataContainer, () => {}, (observer, students) => {observer.next(students[index])},
      (students, studentContainer) => studentContainer.data = students, () => {},
      url, fakeData
    );
  }

  private curStudentsContainer = new DataContainer<Student[][], Student[]>();
  private curStudentsObservables = [
    this.makeStudentObservable(0, "/api/teacher/getstudents?current=true", this.curStudentsContainer, CUR_STUDENT_DATA),
    this.makeStudentObservable(1, "/api/teacher/getstudents?current=true", this.curStudentsContainer, CUR_STUDENT_DATA)
  ]

  private nextStudentsContainer = new DataContainer<Student[][], Student[]>();
  private nextStudentsObservables = [
    this.makeStudentObservable(0, "/api/teacher/getstudents?current=false", this.nextStudentsContainer, NEXT_STUDENT_DATA),
    this.makeStudentObservable(1, "/api/teacher/getstudents?current=false", this.nextStudentsContainer, NEXT_STUDENT_DATA)
  ]

  private allStudentsContainer = new DataContainer<Student[], Student[]>();
  private allStudentsObservable = getHttpCacheObervable<Student[], Student[], Student[], Student[]>(
    this.http, this.allStudentsContainer,() => {}, (observer, students) => observer.next(students), 
    (students, studentsContainer) => studentsContainer.data = students, () => {},
    "/api/teacher/getall?current=false", ALL_STUDENT_DATA
  )

  // private curStudentsObservable = (index: number) => new Observable<Student[]>((observer) => {
  //   if (!this.curIsReqested) {
  //     this.curIsReqested = true;
  //     // this.http.get<RawTeacher[]>("/api/student/getteachers?current=true").
  //     of(CUR_STUDENT_DATA).
  //       subscribe(students => {
  //         this.curStudents = students;
  //         observer.next(students[index]);
  //       }, (err) => this.log(err));
  //   } else {
  //     observer.next(this.curStudents[index])
  //   }
  // })
  // private curIsReqested: boolean = false;
  // private curStudents: Student[][] = new Array<Array<Student>>();
  //
  // private nextStudentsObservable = (index: number) => new Observable<Student[]>((observer) => {
  //   this.nextStudentsObservers[index] = observer;
  //   if (this.nextStudents.length == 0) {
  //     // this.http.get<RawTeacher[]>("/api/student/getteachers?current=false").
  //     of(NEXT_STUDENT_DATA).
  //       subscribe((students) => {
  //         this.nextStudents = students;
  //         observer.next(students[index]);
  //       }, (err) => this.log(err));
  //   } else {
  //     observer.next(this.nextStudents[index]);
  //   }
  //   let observers = this.nextStudentsObservers
  //   return {unsubscribe() {observers.splice(observers.indexOf(observer),1)}}
  // });
  // private nextStudentsObservers: Observer<Student[]>[] = new Array<Observer<Student[]>>(2);
  // private nextStudents: Student[][] = new Array<Array<Student>>();

  // private allStudentsObservable = new Observable<Student[]>((observer) => {
  //   var list = this.allStudentsObservers;
  //   list.push(observer)
  //   if (this.allStudents.length == 0) {
  //     // this.http.get<RawTeacher[]>("/api/teacher/getall?current=false").
  //     of(ALL_STUDENT_DATA).
  //       subscribe(students => this.updateAllStudentsObservers(students),
  //       (err) => this.log(err));
  //   } else {
  //     observer.next(this.allStudents)
  //   }
  //   return {unsubscribe() {list.splice(list.indexOf(observer),1)}}
  // });
  // private allStudentsObservers: Observer<Student[]>[] = new Array<Observer<Student[]>>();
  // private allStudents: Student[] = new Array<Student>();

  constructor(private http: HttpClient) { }

  private updateAllStudentsObservers(students: Student[]) {
    this.allStudentsContainer.data = students;
    this.allStudentsContainer.next();
  }

  getCurStudents(): Observable<Student[]>[] {
    return this.curStudentsObservables;
  }

  getNextStudents(): Observable<Student[]>[] {
    return this.nextStudentsObservables;
  }

  getAllStudents(): Observable<Student[]> {
    return this.allStudentsObservable;
  }

  getStudent(email: string): Student {
    return this.allStudentsContainer.data.find(value => value.Email == email)
  }

  addStudent(block: number, email: string) {
    let student = this.getStudent(email);
    this.nextStudentsContainer.data[block].push(student);
    this.nextStudentsContainer.next();
    this.http.post("/api/teacher/addStudent", {Key: student.ID, Block: block})
    .subscribe(() => {}, (err) => this.log(err));
  }

  removeStudent(block: number, email: string) {
    let student = this.nextStudentsContainer.data[block].find(student => student.Email == email);
    this.nextStudentsContainer.data[block].splice(this.nextStudentsContainer.data[block].indexOf(student),1);
    this.nextStudentsContainer.next();
    this.http.post("/api/teacher/removestudent", {Key: student.ID, Block: block}).
      subscribe(() => {}, (err) => this.log(err));
  }

  log(output) {
    console.log(output);
  }

}

const CUR_STUDENT_DATA: Student[][] = [
  [
    {ID: '11', Name: 'Student11', Email: 'student11', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
    {ID: '12', Name: 'Student12', Email: 'student12', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
    {ID: '13', Name: 'Student13', Email: 'student13', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
    {ID: '14', Name: 'Student14', Email: 'student14', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
    {ID: '15', Name: 'Student15', Email: 'student15', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
  ],
  [
    {ID: '16', Name: 'Student16', Email: 'student16', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
    {ID: '17', Name: 'Student17', Email: 'student17', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
    {ID: '18', Name: 'Student18', Email: 'student18', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
    {ID: '19', Name: 'Student19', Email: 'student19', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
    {ID: '20', Name: 'Student20', Email: 'student20', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
  ]
]

const NEXT_STUDENT_DATA: Student[][] = [
  [
    {ID: '1', Name: 'Student1', Email: 'student1', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
    {ID: '2', Name: 'Student2', Email: 'student2', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
    {ID: '3', Name: 'Student3', Email: 'student3', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
    {ID: '4', Name: 'Student4', Email: 'student4', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
    {ID: '5', Name: 'Student5', Email: 'student5', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
  ],
  [
    {ID: '6', Name: 'Student6', Email: 'student6', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
    {ID: '7', Name: 'Student7', Email: 'student7', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
    {ID: '8', Name: 'Student8', Email: 'student8', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
    {ID: '9', Name: 'Student9', Email: 'student9', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
    {ID: '10', Name: 'Student10', Email: 'student10', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
  ]
]

const ALL_STUDENT_DATA = [
  {ID: '1', Name: 'Student1', Email: 'student1', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
  {ID: '2', Name: 'Student2', Email: 'student2', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
  {ID: '3', Name: 'Student3', Email: 'student3', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
  {ID: '4', Name: 'Student4', Email: 'student4', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
  {ID: '5', Name: 'Student5', Email: 'student5', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
  {ID: '6', Name: 'Student6', Email: 'student6', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
  {ID: '7', Name: 'Student7', Email: 'student7', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
  {ID: '8', Name: 'Student8', Email: 'student8', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
  {ID: '9', Name: 'Student9', Email: 'student9', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
  {ID: '10', Name: 'Student10', Email: 'student10', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
  {ID: '11', Name: 'Student11', Email: 'student11', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
  {ID: '12', Name: 'Student12', Email: 'student12', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
  {ID: '13', Name: 'Student13', Email: 'student13', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
  {ID: '14', Name: 'Student14', Email: 'student14', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
  {ID: '15', Name: 'Student15', Email: 'student15', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
  {ID: '16', Name: 'Student16', Email: 'student16', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
  {ID: '17', Name: 'Student17', Email: 'student17', Grade: 9, Teacher1: '', Teacher2: '', Current: false},
  {ID: '18', Name: 'Student18', Email: 'student18', Grade: 10, Teacher1: '', Teacher2: '', Current: false},
  {ID: '19', Name: 'Student19', Email: 'student19', Grade: 11, Teacher1: '', Teacher2: '', Current: false},
  {ID: '20', Name: 'Student20', Email: 'student20', Grade: 12, Teacher1: '', Teacher2: '', Current: false},
]