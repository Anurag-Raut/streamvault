
import { ping } from 'ldrs'
import { motion } from 'framer-motion'
import { bouncy } from 'ldrs'

bouncy.register()



ping.register()

export default function Loader() {
    return (

        <div>
            <l-ping
                size="120"
                speed="2"
                color="purple"

            ></l-ping>
            <motion.div
                initial={{ opacity: 1 }}
                animate={{ opacity: 0.5 }}
                transition={{ repeat: Infinity, duration: 1, repeatType: "mirror" }}
                className='text-xl font-bold flex'
            >
                Loading
                <div className='ml-2'>
                    <l-bouncy
                        size="17"
                        speed="1.75"
                        color="white"

                    ></l-bouncy>
                </div>

            </motion.div>

        </div>
    )


}