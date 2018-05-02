import { NgModule }             from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  { 
    path: 'student', 
    loadChildren: 'app/student-dashboard/student-dashboard.module#StudentDashboardModule'
  },
  {
    path: '',
    redirectTo: 'student',
    pathMatch: 'full'
  }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {}