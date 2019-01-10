import * as firebase from 'firebase/app'
require("firebase/auth");
require("firebase/database");
require("firebase/firestore")
// const ENV = process.env.NODE

// Initialize Firebase
var config =  {
    apiKey: "",
    authDomain: "",
    databaseURL: "",
    projectId: "",
    storageBucket: "",
    messagingSenderId: ""
}
  
export let appFirebase = firebase.initializeApp(config);

// export const messaging = appFirebase.messaging();
export const FirebaseAuth = firebase.auth;
export const FirebaseDatabase = firebase.database;
// export const firebaseStorage = firebase.storage;
export const FacebookProvider = new firebase.auth.FacebookAuthProvider();
export const FirebaseFirestore = firebase.firestore
export default firebase;
