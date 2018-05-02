import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ChangeDialogComponent } from './change-dialog.component';

describe('ChangeDialogComponent', () => {
  let component: ChangeDialogComponent;
  let fixture: ComponentFixture<ChangeDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ChangeDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ChangeDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
