import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { StudentDashboardRoutingModule } from './student-dashboard-routing.module';
import { StudentDashboardComponent } from './student-dashboard/student-dashboard.component';
import { MaterialImportModule } from '../material-import.module';
import { TeacherService } from '../teacher.service';
import { ChangeDialogComponent } from './change-dialog/change-dialog.component';

@NgModule({
  imports: [
    CommonModule,
    StudentDashboardRoutingModule,
    MaterialImportModule
  ],
  entryComponents: [
    ChangeDialogComponent
  ],
  declarations: [StudentDashboardComponent, ChangeDialogComponent]
})
export class StudentDashboardModule { }
