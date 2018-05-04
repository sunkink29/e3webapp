import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs/observable/of';
import { Observer } from 'rxjs/Observer';
import { HttpClient } from '@angular/common/http';
import 'rxjs/add/operator/map'

import { Teacher, Block } from './Interfaces'

interface RawTeacher {
	ID: string
	Email, Name: string
  Block1 : Block
  Block2 : Block
	Current: boolean
}

@Injectable()
export class TeacherService {

  curClassesObservable = new Observable<Teacher[]>((observer) => {
    if (this.curClasses.length == 0) {
      this.http.get<RawTeacher[]>("/api/student/getteachers?current=true").
        map(value => this.mapRawTeacherToTeacher(value)).
        subscribe((teachers) => {
          this.curClasses = teachers;
          observer.next(this.curClasses);
        }, (err) => this.log(err));
    } else {
      observer.next(this.curClasses);
    }
  });
  curClasses: Teacher[] = new Array<Teacher>();

  nextClassesObservable = new Observable<Teacher[]>((observer) => {
    this.nextClassesObserver = observer;
    if (this.nextClasses.length == 0) {
      this.http.get<RawTeacher[]>("/api/student/getteachers?current=false").
        map(value => this.mapRawTeacherToTeacher(value)).
        subscribe((teachers) => {
          this.nextClasses = teachers;
          observer.next(this.nextClasses);
        }, (err) => this.log(err));
    } else {
      observer.next(this.nextClasses);
    }
    return {unsubscribe() {this.nextClassesObserver = null}}
  });
  nextClassesObserver: Observer<Teacher[]>;
  nextClasses: Teacher[] = new Array<Teacher>();

  allClassesObservable = new Observable<Teacher[]>((observer) => {
    var list = this.allClassesObservers;
    list.push(observer)
    if (this.allClasses.length == 0) {
      this.http.get<RawTeacher[]>("/api/teacher/getall?current=false").
        map(value => this.mapRawTeacherToTeacher(value)).
        subscribe(teachers => this.updateAllClassesObservers(teachers),
        (err) => this.log(err));
    } else {
      observer.next(this.allClasses)
    }
    return {unsubscribe() {list.splice(list.indexOf(observer),1)}}
  });
  allClassesObservers: Observer<Teacher[]>[] = new Array<Observer<Teacher[]>>();
  allClasses: Teacher[] = new Array<Teacher>();

  constructor(private http: HttpClient) { }

  private updateAllClassesObservers(teachers: Teacher[]) {
    this.allClasses = teachers;
    this.allClassesObservers.forEach(observer => observer.next(teachers));
  }

  private mapRawTeacherToTeacher(teachers: RawTeacher[]): Teacher[] {
    let newTeachers: Teacher[] = new Array<Teacher>();
    teachers.forEach( value => {
      let teacher: Teacher;
      if (value != null) {
        teacher = {ID: value.ID, Email: value.Email, Name: value.Name,
          Blocks: [value.Block1, value.Block2], Current: value.Current};
      } else {
        teacher = {ID: "", Email: "", Name: "", Blocks: [
          {Subject: "", Description: "", CurSize: 0, MaxSize: 1, RoomNumber: 0, BlockOpen: true},
          {Subject: "", Description: "", CurSize: 0, MaxSize: 1, RoomNumber: 0, BlockOpen: true}
        ], Current: false}
      }
      newTeachers.push(teacher);
    });
    return newTeachers;
  }

  getCurClasses(): Observable<Teacher[]> {
    return this.curClassesObservable;
  }

  getNextClasses(): Observable<Teacher[]> {
    return this.nextClassesObservable;
  }

  nextContainsClass(Email: string): boolean {
    return this.nextClasses.findIndex(value => value.Email == Email) != -1;
  }

  getAllClasses(): Observable<Teacher[]> {
    return this.allClassesObservable;
  }

  getClass(Email: string): Teacher {
    return this.allClasses.find(value => value.Email == Email)
  }

  setStudentClass(Email: string, block: number) {
    let last = this.getClass(this.nextClasses[block].Email)
    if (last != undefined) {
      last.Blocks[block].CurSize--;
    }
    let next = this.getClass(Email);
    next.Blocks[block].CurSize++;
    this.nextClasses[block] = next;
    this.nextClassesObserver.next(this.nextClasses);
    this.http.post("/api/student/setteacher", {ID: next.ID, Block: block}).
      subscribe(() => {}, (err) => this.log(err));
  }

  log(output) {
    console.log(output);
  }

}

const CUR_TEACHER_DATA: Teacher[] = [
  {ID: "1", Email: 'email1', Name: "teacher1", Blocks: [
    {Subject: "subject1", Description: "description1", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    {Subject: "subject2", Description: "description2", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  ], Current: true},
  {ID: "2", Email: 'email2', Name: "teacher2", Blocks: [
    {Subject: "subject3", Description: "description3", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    {Subject: "subject4", Description: "description4", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: true}
]

const NEXT_TEACHER_DATA: Teacher[] = [
  {ID: "3", Email: 'email3', Name: "teacher3", Blocks: [
    {Subject: "subject5", Description: "description5", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    {Subject: "subject6", Description: "description6", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  ], Current: false},
  {ID: "4", Email: 'email4', Name: "teacher4", Blocks: [
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: false}
]

const All_TEACHER_DATA: Teacher[] = [
  {ID: "1", Email: 'email1', Name: "teacher1", Blocks: [
    {Subject: "subject1", Description: "description1", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    {Subject: "subject2", Description: "description2", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  ], Current: true},
  {ID: "2", Email: 'email2', Name: "teacher2", Blocks: [
    {Subject: "subject3", Description: "description3", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    {Subject: "subject4", Description: "description4", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: true},
  {ID: "3", Email: 'email3', Name: "teacher3", Blocks: [
    {Subject: "subject5", Description: "description5", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    {Subject: "subject6", Description: "description6", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  ], Current: false},
  {ID: "4", Email: 'email4', Name: "teacher4", Blocks: [
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: false},
  {ID: "5", Email: 'email5', Name: "teacher5", Blocks: [
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: false},
  {ID: "6", Email: 'email6', Name: "teacher6", Blocks: [
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: false},
  {ID: "7", Email: 'email7', Name: "teacher7", Blocks: [
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: false},
  
]