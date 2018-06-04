import { Observer, OperatorFunction, Observable, of } from "rxjs";
import { map } from "rxjs/operators";
import { HttpClient } from '@angular/common/http';

export class DataObserver<D ,T> {
  observer: Observer<T>
  next: (observer: Observer<T>, value: D) => void
  constructor(observer: Observer<T>, next: (observer: Observer<T>, value: D) => void) {this.observer = observer; this.next = next}
}

export class DataContainer<D, T> {
  data: D
  requestSent: boolean = false
  observers: DataObserver<D, T>[] = new Array();
  next() {this.observers.forEach(observer => observer.next(observer.observer, this.data))}
}

export function getHttpCacheObervable<R, T, D, I>(
  http: HttpClient,
  dataContainer: DataContainer<D, T>,
  onSubscribe: (observer: Observer<T>, data: D) => void,
  onNext: (observer: Observer<T>, data: D) => void,
  onRequest: (data: I, dataContainer: DataContainer<D, T>) => void,
  onUnsubscribe: (observer: Observer<T>) => void,
  url: string,fakeRequestData?: R, pipe?: OperatorFunction<R, I>
): Observable<T> {
  return new Observable<T>((observer) => {
    dataContainer.observers.push(new DataObserver<D, T>(observer, onNext))
    if (!dataContainer.requestSent) {
      if (pipe == null) {pipe = map<any,any>(value => value)}
      // http.get<R>(url).
      of(fakeRequestData).
      pipe(pipe).subscribe(data => {
        onRequest(data, dataContainer);
        dataContainer.observers.forEach(observer => onSubscribe(observer.observer, dataContainer.data));
        onNext(observer, dataContainer.data);
      }, (err) => this.log(err));
    } else if (dataContainer.data != null) {
      onSubscribe(observer, dataContainer.data);
      onNext(observer, dataContainer.data);
    }
    return {unsubscribe() {onUnsubscribe(observer); 
      dataContainer.observers.splice(dataContainer.observers.map(observer => observer.observer).indexOf(observer),1)}};
  })
}

function log(output) {
  console.log(output);
}