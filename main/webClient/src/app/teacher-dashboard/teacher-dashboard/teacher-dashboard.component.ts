import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { ActivatedRoute, Router, ParamMap } from '@angular/router';
import { MatDialog, MatDialogRef } from '@angular/material';
import { map, filter } from 'rxjs/operators';

import { Student, Block } from '../../Interfaces';
import { StudentService } from '../../student.service';
import { EditDialogComponent } from '../edit-dialog/edit-dialog.component';
import { TeacherService } from '../../teacher.service';
import { AddDialogComponent } from '../add-dialog/add-dialog.component';


@Component({
  selector: 'app-teacher-dashboard',
  templateUrl: './teacher-dashboard.component.html',
  styleUrls: ['./teacher-dashboard.component.css']
})
export class TeacherDashboardComponent implements OnInit {

  nextDisplayedColumns = ['student', 'button'];
  curDisplayedColumns = ['student'];

  nextStudents$: Observable<Student[]>[] = this.studentService.getNextStudents();
  curStudents$: Observable<Student[]>[] = this.studentService.getCurStudents();
  // nextStudents: Student[][] = NEXT_STUDENT_DATA;
  // curStudents: Student[][] = CUR_STUDENT_DATA;

  blockID: number;

  constructor(
    private studentService: StudentService,
    private teacherService: TeacherService,
    private route: ActivatedRoute,
    private router: Router,
    public dialog: MatDialog
  ) { }

  ngOnInit() {
    let editDialogRef: MatDialogRef<EditDialogComponent, Block>;
    let addDialogRef: MatDialogRef<AddDialogComponent, string>;
    setTimeout(() => {
      this.route.url.subscribe(url => {
        const editIndex = url.findIndex((value) => value.toString() === 'edit');
        const addIndex = url.findIndex(value => value.toString() === 'add');
        if (editIndex !== -1) {
          editDialogRef = this.dialog.open(EditDialogComponent, {
            data: {blockID: this.blockID}
          });
          editDialogRef.afterClosed().subscribe((block) => {
            if (block != null) {
              this.teacherService.setNextBlock(this.blockID, block);
            }
            this.router.navigate(['/teacher']);
          });
        }
        if (addIndex != -1) {
          addDialogRef = this.dialog.open(AddDialogComponent, {
            data: {blockID: this.blockID}
          });
          addDialogRef.afterClosed().subscribe(email => {
            if (email != null) {
              this.studentService.addStudent(this.blockID, email);
            }
            this.router.navigate(['/teacher']);
          })
        }
      });
    }, );
    this.route.paramMap.pipe(map((params: ParamMap) => params.get('id'))).subscribe(
      id => {this.blockID = +id; }, () => {});
  }

  removeStudent(email: string, block: number) {
    this.studentService.removeStudent(block, email);
  }

}

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
];

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
];
