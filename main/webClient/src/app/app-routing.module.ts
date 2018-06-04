import { NgModule }             from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  { 
    path: 'student', 
    loadChildren: 'app/student-dashboard/student-dashboard.module#StudentDashboardModule'
  },
  {
    path: 'teacher', 
    loadChildren: 'app/teacher-dashboard/teacher-dashboard.module#TeacherDashboardModule'
  },
  {
    path: '',
    redirectTo: 'teacher',
    pathMatch: 'full'
  }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {}