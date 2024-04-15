"use client"
import { Fragment } from 'react'
import { Menu, Transition } from '@headlessui/react'
import { RiVideoAddFill } from 'react-icons/ri'

function classNames(...classes:any) {
  return classes.filter(Boolean).join(' ')
}

export default function Example() {
  return (
    <Menu as="div" className="relative inline-block text-left z-50 ">
      <div>
        <Menu.Button className="inline-flex w-full justify-center gap-x-1.5 rounded-md bg-background3 px-3 py-2 text-sm font-semibold text-white shadow-sm  hover:bg-card">
      
          <RiVideoAddFill className="fill-primary" size={28} />        </Menu.Button>
      </div>

      <Transition
        as={Fragment}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
        
      >
        <Menu.Items className="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-background3 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none ">
          <div className="py-1">
            <Menu.Item>
              {({ active }) => (
                <a
                  href="/dashboard"
                  className={classNames(
                    active ? 'bg-background4 ' : '',
                    'block px-4 py-2 text-sm text-white z-10'
                  )}
                >
                  Start stream
                </a>
              )}
            </Menu.Item>
            <Menu.Item>
              {({ active }) => (
                <a
                  href="/upload"
                  className={classNames(
                    active ? 'bg-background4 ' : '',
                    'block px-4 py-2 text-sm text-white '
                  )}
                >
                  Upload Videos
                </a>
              )}
            </Menu.Item>

      
          </div>
        </Menu.Items>
      </Transition>
    </Menu>
  )
}
