import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs/observable/of';

import { Teacher } from './Interfaces'
import { Observer } from 'rxjs/Observer';

@Injectable()
export class TeacherService {

  nextClassesObservable = new Observable<Teacher[]>((observer) => {
    this.nextClassesObserver = observer;
    observer.next(NEXT_TEACHER_DATA);
    return {unsubscribe() {this.nextClassesObserver = null}}
  })
  nextClassesObserver: Observer<Teacher[]>;

  allClassesObservers: Observer<Teacher[]>[] = new Array<Observer<Teacher[]>>();

  allClassesObservable = new Observable<Teacher[]>((observer) => {
    var list = this.allClassesObservers;
    list.push(observer)
    observer.next(All_TEACHER_DATA)
    return {unsubscribe() {list.splice(list.indexOf(observer))}}
  })

  constructor() { }

  getCurClasses(): Observable<Teacher[]> {
    return of(CUR_TEACHER_DATA);
  }

  getNextClasses(): Observable<Teacher[]> {
    return this.nextClassesObservable;
  }

  getAllClasses(): Observable<Teacher[]> {
    return this.allClassesObservable;
  }

}

const CUR_TEACHER_DATA: Teacher[] = [
  {ID: "1", Email: 'email1', Name: "teacher1", Blocks: [
    {Subject: "subject1", Description: "description1", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
  ], Current: true},
  {ID: "2", Email: 'email2', Name: "teacher2", Blocks: [
    {Subject: "subject4", Description: "description4", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: false}
  ], Current: true}
]

const NEXT_TEACHER_DATA: Teacher[] = [
  {ID: "3", Email: 'email3', Name: "teacher3", Blocks: [
    {Subject: "subject5", Description: "description5", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
  ], Current: false},
  {ID: "4", Email: 'email4', Name: "teacher4", Blocks: [
    {Subject: "subject8", Description: "description8", CurSize: 0, MaxSize: 5, RoomNumber: 567, BlockOpen: true}
  ], Current: false}
]

const All_TEACHER_DATA: Teacher[] = [
  {ID: "1", Email: 'email1', Name: "teacher1", Blocks: [
    {Subject: "subject1", Description: "description1", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    {Subject: "subject2", Description: "description2", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  ], Current: true},
  {ID: "2", Email: 'email2', Name: "teacher2", Blocks: [
    {Subject: "subject3", Description: "description3", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: false},
    {Subject: "subject4", Description: "description4", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: false}
  ], Current: true},
  {ID: "3", Email: 'email3', Name: "teacher3", Blocks: [
    {Subject: "subject5", Description: "description5", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true},
    {Subject: "subject6", Description: "description6", CurSize: 0, MaxSize: 10, RoomNumber: 123, BlockOpen: true}
  ], Current: false},
  {ID: "4", Email: 'email4', Name: "teacher4", Blocks: [
    {Subject: "subject7", Description: "description7", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: false},
    {Subject: "subject8", Description: "description8", CurSize: 5, MaxSize: 5, RoomNumber: 567, BlockOpen: false}
  ], Current: false}
]