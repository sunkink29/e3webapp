import { Injectable } from '@angular/core';
import { Observable ,  of ,  Observer ,  Subscription, OperatorFunction, observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';

import { Teacher, Block } from './Interfaces'
import { DataContainer, getHttpCacheObervable } from './HttpCacheObervable'

class RawTeacher {
	ID: string
  Email: string
  Name: string
  Block1 : Block
  Block2 : Block
	Current: boolean
}

@Injectable()
export class TeacherService {

  private curClassesObservable = getHttpCacheObervable<RawTeacher[], Teacher[], Teacher[],Teacher[]>(
    this.http, new DataContainer<Teacher[], Teacher[]>(), () => {},
    (observer, teachers) => observer.next(teachers),
    (teachers, dataContainer) => dataContainer.data = teachers,
    () => {}, "/api/student/getteachers?current=true", CUR_TEACHER_DATA, 
    map(value => this.mapRawTeacherToTeacher(value))
  )

  private getNextFuncs() {
    let subscription: Subscription;
    const classRef = this;
    function updateClasses(observer, data: string[]) {
      if (subscription) {
        subscription.unsubscribe();
      }
      subscription = classRef.getAllClasses().subscribe(() => {
        let output = data.map(email => classRef.getClass(email));
        observer.next(output);
      });
    }
    return {updateClasses, unsubscribe() {subscription.unsubscribe()}}
  }
  private nextFuncs = this.getNextFuncs();
  private nextClassesContainer = new DataContainer<string[], Teacher[]>();
  private nextClassesObservable = getHttpCacheObervable<RawTeacher[], Teacher[], string[],Teacher[]>(
    this.http, this.nextClassesContainer, 
    (observer, teachers) => this.nextFuncs.updateClasses(observer, teachers), () => {},
    (teachers, dataContainer) => {
      dataContainer.observers.forEach(observer => observer.observer.next(teachers));
      dataContainer.data = teachers.map(value => value.Email)}, 
    () => this.nextFuncs.unsubscribe(), "/api/student/getteachers?current=false",
    NEXT_TEACHER_DATA, map(value => this.mapRawTeacherToTeacher(value))
  )

  private allClassesContainer = new DataContainer<Teacher[], Teacher[]>();
  private allClassesObservable = getHttpCacheObervable<RawTeacher[], Teacher[], Teacher[],Teacher[]>(
    this.http, this.allClassesContainer, () => {},
    (observer, teachers) => observer.next(teachers),
    (teachers, dataContainer) => dataContainer.data = teachers, 
    () => {},"/api/teacher/getall?current=false", All_TEACHER_DATA, 
    map(value => this.mapRawTeacherToTeacher(value))
  )

  private blockContainer = new DataContainer<Block[], Block>();
  private makeNextBlockObservable(index: number): Observable<Block> {
    return getHttpCacheObervable<Block[],Block,Block[], Block[]>(
      this.http, this.blockContainer, () => {},
      (observer, blocks) => observer.next(blocks[index]), 
      (blocks, blockContainer) => blockContainer.data = blocks,
      () => {}, "/api/teacher/getblocks", NEXT_BLOCK_DATA
    )
  }
  private nextBlockObservable = [
    this.makeNextBlockObservable(0),
    this.makeNextBlockObservable(1)
  ]

  // private curClassesObservableOld = new Observable<Teacher[]>((observer) => {
  //   if (this.curClasses.length == 0) {
  //     // this.http.get<RawTeacher[]>("/api/student/getteachers?current=true").
  //     of(CUR_TEACHER_DATA).pipe(
  //       map(value => this.mapRawTeacherToTeacher(value))
  //     ).subscribe(teachers => {
  //         this.curClasses = teachers;
  //         observer.next(teachers);
  //       }, (err) => this.log(err));
  //   } else {
  //     observer.next(this.curClasses)
  //   }
  // });
  // private curClasses: Teacher[] = new Array<Teacher>();

  // private nextClassesObservableOld = new Observable<Teacher[]>((observer) => {
  //   let subscription: Subscription;
  //   const classRef = this;
  //   function updateClasses() {
  //     subscription = classRef.getAllClasses().subscribe(() => {
  //       let output = classRef.nextClasses.map(email => classRef.getClass(email));
  //       observer.next(output);
  //     });
  //   }

  //   if (this.nextClasses.length == 0) {
  //     // this.http.get<RawTeacher[]>("/api/student/getteachers?current=false").
  //     of(NEXT_TEACHER_DATA).pipe(
  //       map(value => this.mapRawTeacherToTeacher(value))
  //     ).subscribe((teachers) => {
  //         this.nextClasses = teachers.map(value => value.Email);
  //         observer.next(teachers);
  //         updateClasses();
  //       }, (err) => this.log(err));
  //   } else {
  //     updateClasses();
  //   }
  //   return {unsubscribe() {subscription.unsubscribe()}}
  // });
  // private nextClasses: string[] = new Array<string>();

  // private allClassesObservableOld = new Observable<Teacher[]>((observer) => {
  //   var list = this.allClassesObservers;
  //   list.push(observer)
  //   if (this.allClasses.length == 0) {
  //     // this.http.get<RawTeacher[]>("/api/teacher/getall?current=false").
  //     of(All_TEACHER_DATA).pipe(
  //       map(value => this.mapRawTeacherToTeacher(value))
  //     ).subscribe(teachers => this.updateAllClassesObservers(teachers),
  //       (err) => this.log(err));
  //   } else {
  //     observer.next(this.allClasses)
  //   }
  //   return {unsubscribe() {list.splice(list.indexOf(observer),1)}}
  // });
  // private allClasses: Teacher[] = new Array<Teacher>();
  // private allClassesObservers: Observer<Teacher[]>[] = new Array<Observer<Teacher[]>>();

  constructor(private http: HttpClient) { }

  private updateAllClassesObservers(teachers: Teacher[]) {
    this.allClassesContainer.data = teachers;
    this.allClassesContainer.next();
  }

  mapRawTeacherToTeacher(teachers: RawTeacher[]): Teacher[] {
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

  getAllClasses(): Observable<Teacher[]> {
    return this.allClassesObservable;
  }

  getClass(Email: string): Teacher {
    return this.allClassesContainer.data.find(value => value.Email == Email)
  }

  setStudentClass(email: string, block: number) {
    if (block < 0 || block >= this.nextClassesContainer.data.length) {return}

    let last = this.getClass(this.nextClassesContainer.data[block]);
    if (last != undefined) {
      last.Blocks[block].CurSize--;
    }
    let next = this.getClass(email);
    next.Blocks[block].CurSize++;
    this.nextClassesContainer.data[block] = email;
    this.updateAllClassesObservers(this.allClassesContainer.data);
    this.http.post("/api/student/setteacher", {ID: next.ID, Block: block}).
      subscribe(() => {}, (err) => this.log(err));
  }

  getNextBlock(blockID: number): Observable<Block> {
    return this.nextBlockObservable[blockID];
  }
  
  setNextBlock(blockID: number, block: Block) {
    this.blockContainer.data[blockID] = block;
    this.blockContainer.next();
  }

  log(output) {
    console.log(output);
  }

}

const CUR_TEACHER_DATA: RawTeacher[] = [
  {ID: "1", Email: 'email1', Name: "teacher1", Block1:
    {Subject: "subject1", Description: "description1", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    Block2:
    {Subject: "subject2", Description: "description2", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  , Current: true},
  {ID: "2", Email: 'email2', Name: "teacher2", Block1:
    {Subject: "subject3", Description: "description3", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    Block2:
    {Subject: "subject4", Description: "description4", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  , Current: true}
]

const NEXT_TEACHER_DATA: RawTeacher[] = [
  {ID: "3", Email: 'email3', Name: "teacher3", Block1:
    {Subject: "subject5", Description: "description5", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    Block2:
    {Subject: "subject6", Description: "description6", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  , Current: false},
  {ID: "4", Email: 'email4', Name: "teacher4", Block1:
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    Block2:
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  , Current: false}
]

const All_TEACHER_DATA: RawTeacher[] = [
  {ID: "1", Email: 'email1', Name: "teacher1", Block1:
    {Subject: "subject1", Description: "description1", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    Block2:
    {Subject: "subject2", Description: "description2", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  , Current: true},
  {ID: "2", Email: 'email2', Name: "teacher2", Block1:
    {Subject: "subject3", Description: "description3", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    Block2:
    {Subject: "subject4", Description: "description4", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  , Current: true},
  {ID: "3", Email: 'email3', Name: "teacher3", Block1:
    {Subject: "subject5", Description: "description5", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    Block2:
    {Subject: "subject6", Description: "description6", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  , Current: false},
  {ID: "4", Email: 'email4', Name: "teacher4", Block1:
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    Block2:
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  , Current: false},
  {ID: "5", Email: 'email5', Name: "teacher5", Block1:
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    Block2:
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  , Current: false},
  {ID: "6", Email: 'email6', Name: "teacher6", Block1:
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    Block2:
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  , Current: false},
  {ID: "7", Email: 'email7', Name: "teacher7", Block1:
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true},
    Block2:
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  , Current: false},
]

const NEXT_BLOCK_DATA: Block[] = [
  {Subject: "subject1", Description: "description1", CurSize: 1, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
  {Subject: "subject2", Description: "description2", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
]