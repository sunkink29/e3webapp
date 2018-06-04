import { Component, OnInit, Inject } from '@angular/core';
import {MAT_DIALOG_DATA} from '@angular/material';
import { Observable } from 'rxjs';

import { Teacher } from '../../Interfaces';
import { TeacherService } from '../../teacher.service';

@Component({
  selector: 'app-change-dialog',
  templateUrl: './change-dialog.component.html',
  styleUrls: ['./change-dialog.component.css']
})
export class ChangeDialogComponent implements OnInit {

  displayedColumns = ['name', 'roomNumber', 'status', 'subject', 'description', 'size', 'changeButton'];

  allTeachers$: Observable<Teacher[]>;

  constructor(
    private teacherService: TeacherService,
    @Inject(MAT_DIALOG_DATA) public blockID: number) { }

  ngOnInit() {
    this.allTeachers$ = this.teacherService.getAllClasses()
  }
}
