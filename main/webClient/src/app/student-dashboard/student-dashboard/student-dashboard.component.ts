import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { ActivatedRoute, Router, ParamMap } from '@angular/router';
import { MatDialog, MatDialogRef } from '@angular/material';
import 'rxjs/add/operator/map';


import { Teacher, Block} from '../../Interfaces';
import { TeacherService } from '../../teacher.service';
import { ChangeDialogComponent } from '../change-dialog/change-dialog.component';

@Component({
  selector: 'app-student-dashboard',
  templateUrl: './student-dashboard.component.html',
  styleUrls: ['./student-dashboard.component.css']
})
export class StudentDashboardComponent implements OnInit {

  indexT = 0
  indexB = 1;
  curDisplayedColumns = ['days', 'name', 'roomNumber', 'status', 'subject', 'description', 'size'];
  nextDisplayedColumns = ['days', 'name', 'roomNumber', 'status', 'subject', 'description', 'size', 'changeButton'];

  curTeachers$: Observable<Teacher[]>;
  nextTeachers$: Observable<Teacher[]>;

  Math = Math;

  blockID: number;

  constructor(
    private teacherService: TeacherService, 
    private route: ActivatedRoute,
    private router: Router,
    public dialog: MatDialog) { }

  ngOnInit() {
    this.curTeachers$ = this.teacherService.getCurClasses();
    this.nextTeachers$ = this.teacherService.getNextClasses();
    let dialogRef;
    setTimeout(() => {this.route.url
      .subscribe(url => {
        let index = url.findIndex((value) => value.toString() == 'change')
        if (index != -1) {
          dialogRef = this.dialog.open(ChangeDialogComponent, {
            height: '30rem',
            width: '80rem',
            data: this.blockID
          });
          dialogRef.afterClosed().subscribe(teacherID => {
            console.log(teacherID);
            this.router.navigate(['/student']);
          });
        }
      });
    },);
    this.route.paramMap.map((params: ParamMap) => params.get('id')).subscribe(
      id => {this.blockID = +id}, () => {})
  }

  getIndexT() {
    this.indexT++;
    return this.indexT % 2;
  }

  getIndexB() {
    this.indexB++;
    return this.indexB % 4;
  }

}