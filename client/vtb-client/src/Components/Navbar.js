import React from 'react'
import { NavLink } from 'react-router-dom'
import { Fragment } from 'react';

export default function Navbar() {
  const userAuthenticated = false;

  return (
    <nav className="navbar navbar-expand-lg navbar-light bg-light">
      <img className="mr-4" src="https://upload.wikimedia.org/wikipedia/commons/7/7c/VTB_Logo_2018.svg" width="60" alt="Kiwi standing on oval" />
      <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
        <span className="navbar-toggler-icon"></span>
      </button>

      <div className="collapse navbar-collapse" id="navbarNav">
        <ul className="navbar-nav">
          <li className="nav-item active">
            <NavLink className="nav-link text-primary" to="/">Главная</NavLink>
          </li>
          {userAuthenticated
            ? (
              <li className="nav-item">
                <NavLink className="nav-link text-primary" to="/">Выйти</NavLink>
              </li>
            ) : (
              <>
                <li className="nav-item">
                  <NavLink className="nav-link text-primary" to="/login">Войти</NavLink>
                </li>
                <li className="nav-item">
                  <NavLink className="nav-link text-primary" to="/signup">Регистрация</NavLink>
                </li>
              </>
            )
          }
        </ul>
      </div>
    </nav>
  )
}
