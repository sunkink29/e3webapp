import { Component, OnInit, Inject } from '@angular/core';
import { Observable } from 'rxjs';
import { Student } from '../../Interfaces';
import { TeacherService } from '../../teacher.service';
import { MAT_DIALOG_DATA } from '@angular/material';
import { StudentService } from '../../student.service';

@Component({
  selector: 'app-add-dialog',
  templateUrl: './add-dialog.component.html',
  styleUrls: ['./add-dialog.component.css']
})
export class AddDialogComponent implements OnInit {

  students$: Observable<Student[]> = this.studentService.getAllStudents();

  constructor(
    private studentService: StudentService,
    @Inject(MAT_DIALOG_DATA) public data: {
      blockID: number}
  ) { }

  ngOnInit() {
  }

}
